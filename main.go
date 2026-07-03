package main

import (
	"fmt"

	"onah-go/config"
	"onah-go/database"
	"onah-go/whatsapp"

	// Import plugins to trigger their init() functions and register commands
	_ "onah-go/plugins"
)

func main() {
	// ASCII Art banner for a premium console experience
	fmt.Println(`
\033[1;36m========================================================
  __  __               _                       
 |  \/  | __ _ _ __   | | __ ___  _ __   __ _  
 | |\/| |/ _' | '__|  | |/ // _ \| '_ \ / _' | 
 | |  | | (_| | |     |   <| (_) | | | | (_| | 
 |_|  |_|\__,_|_|     |_|\_\\___/|_| |_|\__,_| \033[1;32mv2.0-Go\033[0m
                                               
     ⚡ High Performance Go WhatsApp Bot ⚡
========================================================\033[0m
	`)

	fmt.Println("[⚙️ ONAH-GO] Loading configuration...")
	err := config.LoadConfig("config.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	fmt.Println("[🗄️ ONAH-GO] Initializing SQLite database...")
	err = database.InitDB(config.ActiveConfig.DatabasePath)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize database: %v", err))
	}

	fmt.Println("[💬 ONAH-GO] Connecting to WhatsApp...")
	whatsapp.InitClient()
	whatsapp.Connect()

	fmt.Println("[🚀 ONAH-GO] Bot is running! Press Ctrl+C to stop.")
	
	// Keep the main thread alive indefinitely.
	// Graceful shutdown is managed by OS signals in whatsapp.Connect()
	select {}
}
