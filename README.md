# ⚡ ONAH-GO: High Performance WhatsApp Bot in Golang

Onah-Go adalah bot WhatsApp modular yang ditulis menggunakan bahasa pemrograman **Golang (Go)** dengan pustaka utama **[whatsmeow](https://github.com/tulir/whatsmeow)** dan penyimpanan sesi/data berbasis **SQLite**. 

Proyek ini dirancang khusus untuk deployment efisiensi tinggi pada **Linux VPS** atau **VPS Panel (Pterodactyl)** dengan konsumsi resource yang sangat minim.

---

## 🚀 Mengapa Memilih Go (whatsmeow) dibanding Node.js (Baileys)?

* **Super Hemat RAM**: Bot NodeJS Baileys umumnya memakan **300MB - 600MB RAM**. Onah-Go hanya memerlukan sekitar **20MB - 50MB RAM**! Sangat hemat biaya dan enteng di VPS murah spesifikasi terendah sekalipun.
* **Instan & Responsif**: Go dicompile langsung ke bahasa mesin (native binary) sehingga eksekusi command bot jauh lebih cepat dan minim latensi dibanding NodeJS.
* **Tombol Interaktif Modern**: Menggunakan protokol **WhatsApp Interactive Message (Native Flow)** terbaru, sehingga tombol/list menu akan merender secara sempurna di Android, iOS, dan WhatsApp Web.

---

## 🌐 Panduan Deploy 24/7 di Linux VPS (Ubuntu/Debian)

Ini adalah metode deployment standar yang paling sering digunakan untuk bot WhatsApp.

### 1. Persiapan Environment
Instal Golang dan build tools (GCC) untuk kebutuhan database SQLite (CGO):
```bash
sudo apt update
sudo apt install build-essential golang -y
```

### 2. Jalankan Bot
Clone repositori, unduh dependensi, lalu jalankan program:
```bash
go mod tidy
go run main.go
```
*Saat pertama kali dijalankan, bot akan menanyakan metode login di terminal: ketik **`1`** untuk pairing code, atau ketik **`2`** untuk QR code.*

### 3. Deploy Menggunakan PM2 (Latar Belakang 24/7)
Kompilasi program ke bentuk binary executable terlebih dahulu:
```bash
# Build binary
go build -o onah-bot main.go

# Jalankan menggunakan PM2
pm2 start ./onah-bot --name onah-go

# Perintah PM2 lainnya
pm2 logs onah-go     # Melihat log bot
pm2 restart onah-go  # Restart bot
pm2 stop onah-go     # Matikan bot
```

---

## 🎛️ Panduan Deploy di VPS Panel (Pterodactyl)

Jika Anda menggunakan Pterodactyl Panel untuk menjalankan bot:

1. Pastikan Server menggunakan **Go Egg** (atau Docker Image yang mendukung Golang/Debian).
2. Upload seluruh source code bot ke File Manager panel (kecuali folder `storage` atau `database.db` sesi Anda).
3. Jika Panel tidak memiliki compiler GCC/CGO, Anda disarankan untuk melakukan **Cross-Compilation** di komputer lokal Anda terlebih dahulu (lihat bagian di bawah), lalu cukup upload file binary hasil build (`onah-bot`) ke Pterodactyl.
4. Set **Startup Command** di panel Anda menjadi:
   ```bash
   ./onah-bot
   ```

---

## 🔀 Panduan Cross-Compilation (Build Linux dari Windows)
Jika Anda mengedit/mengembangkan bot di komputer Windows tetapi ingin men-deploy-nya ke VPS Linux atau Panel, Anda bisa melakukan cross-compile dengan perintah ini di terminal Windows Anda:

* **Di CMD**:
  ```cmd
  set GOOS=linux
  set GOARCH=amd64
  set CGO_ENABLED=1
  go build -o onah-bot main.go
  ```
* **Di PowerShell**:
  ```powershell
  $env:GOOS="linux"
  $env:GOARCH="amd64"
  $env:CGO_ENABLED="1"
  go build -o onah-bot main.go
  ```
Setelah itu, file `onah-bot` (tanpa akhiran `.exe`) akan terbentuk. Anda tinggal meng-upload file tersebut ke VPS Linux atau Pterodactyl Panel Anda dan langsung menjalankannya!

---

## 💻 Cara Menjalankan Bot di Windows (Lokal/Pengembangan)

Jika Anda ingin mengetes atau mengembangkan bot di komputer Windows:

1. Unduh dan pasang compiler C portable **[w64devkit](https://github.com/skeeto/w64devkit/releases)** di komputer Anda (misal letakkan di `C:\w64devkit`).
2. Jalankan bot lewat terminal PowerShell dengan perintah:
   ```powershell
   $env:PATH = 'C:\w64devkit\bin;' + $env:PATH; go run main.go
   ```

---

## ⚙️ Konfigurasi (`config.json`)

Sunting file `config.json` di direktori utama untuk melakukan penyesuaian:
```json
{
  "owner_number": "6285802569316",
  "owner_name": "Izan",
  "bot_name": "ONAH-GO",
  "prefixes": [".", "!", "/"],
  "database_path": "storage/database.db",
  "limit_default": 20,
  "pairing_code_enabled": true,
  "pairing_number": "6285165613515",
  "api_keys": {
    "gemini": "MASUKKAN_API_KEY_GEMINI_ANDA_DI_SINI"
  },
  "messages": {
    "wait": "Onah sedang memproses...",
    "error": "Oops, terjadi kesalahan!",
    "owner_only": "Fitur ini khusus untuk Owner bot!"
  }
}
```
> **⚠️ PENTING**: Jangan pernah mempublikasikan file `config.json` yang berisi API Key asli Anda ke repositori GitHub publik.

---

## 📂 Struktur File Utama
* `main.go`: Bootstrapper utama program.
* `config/`: Kode parsing file `config.json`.
* `database/`: Controller data pengguna (Limit & Balance) menggunakan SQLite.
* `whatsapp/`: Logic penanganan socket `whatsmeow` & interceptor pesan masuk.
* `plugins/`: Tempat menulis modul-modul fitur command bot.
