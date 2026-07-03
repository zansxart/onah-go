# ⚡ Onah-Go: High Performance Go WhatsApp Bot

Onah-Go adalah duplikasi dan porting dari bot WhatsApp **Onah (Node.js/JS)** ke dalam bahasa pemrograman **Golang (Go)** menggunakan pustaka **[whatsmeow](https://github.com/tulir/whatsmeow)** dan penyimpanan database **SQLite**.

Bot ini dirancang agar sangat hemat penggunaan RAM (~20MB - 50MB) dan CPU, menjadikannya sangat cocok untuk dijalankan di VPS dengan spesifikasi rendah sekalipun.

---

## 🛠️ Persyaratan Sistem
Sebelum menjalankan bot, pastikan sistem Anda telah memiliki:
1. **Golang** (Versi 1.21 ke atas)
2. **Kompiler C (GCC)** karena SQLite (`go-sqlite3`) memerlukan CGO untuk kompilasi:
   - **Di Windows**: Gunakan [w64devkit](https://github.com/skeeto/w64devkit/releases/tag/v1.23.0) (Portable GCC) atau MSYS2.
   - **Di Linux VPS**: Cukup install tool build-essential:
     ```bash
     sudo apt update && sudo apt install build-essential -y
     ```

---

## ⚙️ Konfigurasi (`config.json`)
Sunting file `config.json` di folder utama proyek untuk menyesuaikan pengaturan bot:
```json
{
  "owner_number": "6285802569316",
  "owner_name": "Izann",
  "bot_name": "markonah-md",
  "prefixes": [".", "!", "/"],
  "database_path": "storage/database.db",
  "limit_default": 20,
  "pairing_code_enabled": true,
  "pairing_number": "6285165613514",
  "api_keys": {
    "gemini": "AIzaSyCMARX..."
  },
  "messages": {
    "wait": "Onah sedang memproses...",
    "error": "Oops, terjadi kesalahan!",
    "owner_only": "Fitur ini khusus untuk Owner bot!"
  }
}
```

---

## 🚀 Cara Menjalankan Bot

### 1. Unduh Dependensi Proyek
Jalankan perintah ini satu kali setelah instalasi awal untuk men-download semua package luar:
```bash
go mod tidy
```

### 2. Jalankan Mode Pengembangan (Development)
Untuk menjalankan bot secara langsung dari source code:

* **Di Windows (jika GCC diinstal di C:\w64devkit)**:
  ```powershell
  $env:PATH = 'C:\w64devkit\bin;' + $env:PATH; go run main.go
  ```
* **Di Linux VPS / Termux**:
  ```bash
  go run main.go
  ```

---

## 📦 Mengompilasi ke Single Binary (Produksi)
Untuk mengompilasi bot menjadi satu file binary mandiri yang siap dideploy tanpa perlu source code lagi:

* **Kompilasi di Windows (.exe)**:
  ```powershell
  $env:PATH = 'C:\w64devkit\bin;' + $env:PATH; go build -o onah-bot.exe main.go
  ```
* **Kompilasi di Linux VPS**:
  ```bash
  go build -o onah-bot main.go
  ```

---

## 🌐 Panduan Deploy 24/7 di Linux VPS
Untuk menjalankan bot secara terus-menerus (24/7) di VPS, Anda memiliki dua opsi populer:

### Opsi A: Menggunakan PM2 (Paling Sederhana)
PM2 juga mendukung eksekusi file binary hasil compile:
```bash
# Jalankan binary bot menggunakan pm2
pm2 start ./onah-bot --name onah-go

# Melihat log bot
pm2 logs onah-go

# Menghentikan bot
pm2 stop onah-go
```

### Opsi B: Menggunakan Systemd Service (Sistem Linux Bawaan)
Buat file service systemd di Linux:
```bash
sudo nano /etc/systemd/system/onah-bot.service
```

Masukkan template berikut:
```ini
[Unit]
Description=Onah Go WhatsApp Bot
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/onah-go
ExecStart=/root/onah-go/onah-bot
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Aktifkan dan jalankan service:
```bash
sudo systemctl enable onah-bot
sudo systemctl start onah-bot
sudo systemctl status onah-bot
```

---

## 📂 Struktur Kode
* `main.go`: Entrypoint program (bootstrap).
* `config/`: Pengaturan pembaca file `config.json`.
* `database/`: Mengelola SQLite user model (Limit & Saldo).
* `whatsapp/`: Logic koneksi `whatsmeow` & routing pesan masuk.
* `plugins/`: Kumpulan fitur command bot.
  - `main_plugin.go`: Command dasar bot (`.ping`, `.register`, `.limit`, `.menu`).
  - `ai_plugin.go`: Integrasi AI Gemini (`.ai`).
  - `downloader_plugin.go`: Template Downloader (`.tiktok`).
