package plugins

import (
	"fmt"
	"strings"
	"time"

	"onah-go/config"
	"onah-go/database"
)

func init() {
	Register(Command{
		Name: "ping",
		Tags: []string{"main"},
		Help: "Memeriksa status dan latency bot",
		Execute: func(ctx *Context) error {
			ctx.React("🕛")
			latency := time.Since(ctx.Event.Info.Timestamp)
			return ctx.Reply(fmt.Sprintf("Pong! 🏓 Latency: %s", latency.Round(time.Millisecond)))
		},
	})

	Register(Command{
		Name: "register",
		Tags: []string{"main"},
		Help: "Mendaftar sebagai pengguna bot",
		Execute: func(ctx *Context) error {
			if ctx.User.Registered {
				return ctx.Reply("⚠️ Kamu sudah terdaftar sebelumnya!")
			}
			
			name := ctx.Query
			if name == "" {
				name = ctx.PushName
			}
			
			ctx.User.Registered = true
			ctx.User.Name = name
			ctx.User.RegisteredAt = time.Now().Format("2006-01-02 15:04:05")

			err := database.UpdateUser(ctx.User)
			if err != nil {
				return err
			}

			return ctx.Reply(fmt.Sprintf("✅ Registrasi berhasil!\nNama: %s\nTanggal: %s", name, ctx.User.RegisteredAt))
		},
	})

	Register(Command{
		Name: "limit",
		Tags: []string{"main"},
		Help: "Melihat limit dan profil pengguna",
		Execute: runProfile,
	})

	Register(Command{
		Name: "profile",
		Tags: []string{"main"},
		Help: "Melihat limit dan profil pengguna",
		Execute: runProfile,
	})

	Register(Command{
		Name: "menu",
		Tags: []string{"main"},
		Help: "Menampilkan semua daftar fitur bot",
		Execute: func(ctx *Context) error {
			ctx.React("📖")

			allCmds := GetCommands()
			tagMap := make(map[string][]Command)
			for _, cmd := range allCmds {
				for _, tag := range cmd.Tags {
					tagMap[tag] = append(tagMap[tag], cmd)
				}
			}

			totalUsers, _ := database.GetTotalUsers()
			registeredUsers, _ := database.GetRegisteredUsers()

			statusReg := "❌ Belum Terdaftar"
			if ctx.User.Registered {
				statusReg = "✅ Terdaftar"
			}
			statusPrem := "🔹 Free"
			if ctx.User.Premium {
				statusPrem = "⭐ Premium"
			}

			header := fmt.Sprintf(`┏━━━〔 *DASHBOARD* 〕━⬣
┃ ✦ Nama   : %s
┃ ✦ Limit  : %d
┃ ✦ Saldo  : Rp %d
┃ ✦ Status : %s (%s)
┗⬣

┏━━━〔 *INFO BOT* 〕━⬣
┃ ✦ Bot Mode    : Public
┃ ✦ Total User  : %d
┃ ✦ Terdaftar   : %d
┃ ✦ Owner       : @%s
┗⬣

`, ctx.PushName, ctx.User.Limit, ctx.User.Money, statusReg, statusPrem, totalUsers, registeredUsers, config.ActiveConfig.OwnerNumber)

			var menuText strings.Builder
			menuText.WriteString(header)

			// Formatting readmore spacer (equivalent to more.repeat(4001) in JS)
			readMore := string(rune(8206)) + strings.Repeat(string(rune(8206)), 4000)
			menuText.WriteString(readMore)
			menuText.WriteString("📖 *DAFTAR MENU ONAN-GO*\n\n")

			for tag, cmds := range tagMap {
				menuText.WriteString(fmt.Sprintf("┏━━━〔 *MENU %s* 〕━⬣\n", strings.ToUpper(tag)))
				for _, cmd := range cmds {
					flag := "🅕"
					if cmd.Limit {
						flag = "🅛"
					} else if cmd.Premium {
						flag = "🅟"
					} else if cmd.OwnerOnly {
						flag = "🅞"
					}
					menuText.WriteString(fmt.Sprintf("┃ ◕ .%s [%s]\n", cmd.Name, flag))
				}
				menuText.WriteString("┗⬣\n\n")
			}

			menuText.WriteString("ℹ️ _Keterangan:_\n_🅕 = Free, 🅛 = Limit, 🅟 = Premium, 🅞 = Owner Only_")

			return ctx.Reply(menuText.String())
		},
	})
}

func runProfile(ctx *Context) error {
	statusReg := "❌ Belum Terdaftar"
	if ctx.User.Registered {
		statusReg = "✅ Terdaftar"
	}
	statusPrem := "🔹 Free"
	if ctx.User.Premium {
		statusPrem = "⭐ Premium"
	}

	profile := fmt.Sprintf(`┏━━━〔 *USER PROFILE* 〕━⬣
┃ ✦ Nama   : %s
┃ ✦ Nomor  : @%s
┃ ✦ Limit  : %d
┃ ✦ Saldo  : Rp %d
┃ ✦ Status : %s (%s)
┗⬣`, ctx.User.Name, strings.Split(ctx.SenderJID, "@")[0], ctx.User.Limit, ctx.User.Money, statusReg, statusPrem)

	return ctx.Reply(profile)
}
