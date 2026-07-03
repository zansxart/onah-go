package whatsapp

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"onah-go/config"
	"onah-go/database"
	"onah-go/plugins"
)

var Client *whatsmeow.Client
var Log waLog.Logger

func InitClient() {
	Log = waLog.Stdout("WhatsMeow", "INFO", true)

	dbLog := waLog.Stdout("Database", "WARN", true)
	// Open connection to whatsmeow store (using same SQLite db but different tables)
	container, err := sqlstore.New(context.Background(), "sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", config.ActiveConfig.DatabasePath), dbLog)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize session database: %v", err))
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		panic(err)
	}

	Client = whatsmeow.NewClient(deviceStore, Log)
	Client.AddEventHandler(eventHandler)
}

func Connect() {
	if Client.Store.ID == nil {
		// New login
		fmt.Println("\n============================================")
		fmt.Println("🔑 METODE LOGIN WHATSAPP")
		fmt.Println("============================================")
		fmt.Println("1. Pairing Code (Tautkan dengan Nomor Telepon)")
		fmt.Println("2. QR Code (Scan Barcode di Layar)")
		fmt.Println("============================================")

		reader := bufio.NewReader(os.Stdin)
		var choice string
		for {
			fmt.Print("Pilih opsi login Anda (1 atau 2): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				choice = "2" // Default to QR if stdin is closed/errored
				break
			}
			choice = strings.TrimSpace(input)
			if choice == "1" || choice == "2" {
				break
			}
			fmt.Println("⚠️ Pilihan tidak valid. Silakan masukkan angka 1 atau 2.")
		}

		if choice == "1" {
			// Pairing code method
			phone := config.ActiveConfig.PairingNumber
			if phone == "" || phone == "pairing_number" {
				fmt.Print("Masukkan nomor HP WhatsApp bot Anda (contoh: 6285xxx): ")
				input, err := reader.ReadString('\n')
				if err == nil {
					phone = strings.TrimSpace(input)
				}
			}

			err := Client.Connect()
			if err != nil {
				panic(err)
			}

			// Clean number format (remove non-digits)
			phone = strings.ReplaceAll(phone, "+", "")
			phone = strings.ReplaceAll(phone, " ", "")
			phone = strings.ReplaceAll(phone, "-", "")

			// Get pairing code
			code, err := Client.PairPhone(context.Background(), phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
			if err != nil {
				fmt.Printf("Failed to get pairing code: %v\n", err)
				return
			}
			fmt.Printf("\n====================================\n")
			fmt.Printf("🔗 PAIRING CODE FOR %s: \n", phone)
			fmt.Printf("👉 \033[1;32m%s\033[0m 👈\n", formatPairingCode(code))
			fmt.Printf("====================================\n\n")
		} else {
			// QR code method
			qrChan, err := Client.GetQRChannel(context.Background())
			if err != nil {
				panic(err)
			}
			err = Client.Connect()
			if err != nil {
				panic(err)
			}
			go func() {
				for qr := range qrChan {
					if qr.Event == "code" {
						fmt.Println("Scan this QR code to login:")
						qrterminal.GenerateHalfBlock(qr.Code, qrterminal.L, os.Stdout)
					} else {
						fmt.Printf("QR Event: %s\n", qr.Event)
					}
				}
			}()
		}
	} else {
		// Already logged in
		err := Client.Connect()
		if err != nil {
			panic(err)
		}
		fmt.Println("Successfully connected to WhatsApp!")
	}

	// Listen for interrupt to close client cleanly
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		Client.Disconnect()
		fmt.Println("Disconnected. Exiting...")
		os.Exit(0)
	}()
}

func formatPairingCode(code string) string {
	if len(code) == 8 {
		return fmt.Sprintf("%s-%s", code[0:4], code[4:8])
	}
	return code
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		handleMessage(v)
	}
}

func getMessageText(msg *waE2E.Message) string {
	if msg == nil {
		return ""
	}
	if msg.Conversation != nil {
		return *msg.Conversation
	}
	if msg.ExtendedTextMessage != nil && msg.ExtendedTextMessage.Text != nil {
		return *msg.ExtendedTextMessage.Text
	}
	if msg.ImageMessage != nil && msg.ImageMessage.Caption != nil {
		return *msg.ImageMessage.Caption
	}
	if msg.VideoMessage != nil && msg.VideoMessage.Caption != nil {
		return *msg.VideoMessage.Caption
	}
	if msg.DocumentMessage != nil && msg.DocumentMessage.Caption != nil {
		return *msg.DocumentMessage.Caption
	}
	return ""
}

func handleMessage(evt *events.Message) {
	body := getMessageText(evt.Message)
	senderJID := evt.Info.Sender.ToNonAD().String()

	// Cetak pesan masuk ke terminal untuk debugging
	if body != "" {
		fmt.Printf("[📩 PESAN MASUK] Dari: %s (%s) | Isi: %s | IsFromMe: %t\n", senderJID, evt.Info.PushName, body, evt.Info.IsFromMe)
	}

	// Abaikan pesan dari nomor bot sendiri
	if evt.Info.IsFromMe {
		return
	}

	if body == "" {
		return
	}

	chatJID := evt.Info.Chat.String()
	pushName := evt.Info.PushName

	// Detect Prefix
	var matchedPrefix string
	hasPrefix := false
	for _, pref := range config.ActiveConfig.Prefixes {
		if strings.HasPrefix(body, pref) {
			matchedPrefix = pref
			hasPrefix = true
			break
		}
	}

	if !hasPrefix {
		return
	}

	// Parse Command and Arguments
	fullCmd := strings.TrimPrefix(body, matchedPrefix)
	parts := strings.Fields(fullCmd)
	if len(parts) == 0 {
		return
	}

	cmdName := strings.ToLower(parts[0])
	args := parts[1:]
	query := strings.TrimSpace(strings.TrimPrefix(fullCmd, parts[0]))

	// Match command in registry
	var matchedCmd *plugins.Command
	for _, cmd := range plugins.GetCommands() {
		if cmd.Name == cmdName {
			matchedCmd = &cmd
			break
		}
	}

	if matchedCmd == nil {
		return
	}

	// Load User or Create if not exists in local DB
	user, err := database.GetUser(senderJID)
	if err != nil {
		fmt.Printf("Database error fetching user: %v\n", err)
		return
	}
	if user == nil {
		// Auto-register user with default limits
		user, err = database.CreateUser(senderJID, pushName, config.ActiveConfig.LimitDefault)
		if err != nil {
			fmt.Printf("Database error creating user: %v\n", err)
			return
		}
	}

	// Prepare Context
	ctx := &plugins.Context{
		Client:    Client,
		Event:     evt,
		SenderJID: senderJID,
		ChatJID:   chatJID,
		PushName:  pushName,
		Command:   cmdName,
		Args:      args,
		Query:     query,
		User:      user,
	}

	// Check constraints
	isOwner := strings.Contains(senderJID, config.ActiveConfig.OwnerNumber)

	if matchedCmd.OwnerOnly && !isOwner {
		ctx.Reply(config.ActiveConfig.Messages["owner_only"])
		return
	}

	// If command requires limit, check and deduct
	if matchedCmd.Limit && !isOwner && !user.Premium {
		if user.Limit <= 0 {
			ctx.Reply("⚠️ Limit harian Anda telah habis! Hubungi owner untuk menambah limit atau mendaftar Premium.")
			return
		}
		user.Limit--
		err = database.UpdateUser(user)
		if err != nil {
			fmt.Printf("Failed to update user limit: %v\n", err)
		}
	}

	// Execute command
	err = matchedCmd.Execute(ctx)
	if err != nil {
		fmt.Printf("Error executing command %s: %v\n", cmdName, err)
		ctx.Reply(fmt.Sprintf("%s\nDetail: %v", config.ActiveConfig.Messages["error"], err))
	}
}
