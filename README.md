# ⚡ ONAH-GO: High Performance WhatsApp Bot in Golang

Onah-Go adalah bot WhatsApp modular yang ditulis menggunakan bahasa pemrograman **Golang (Go)** dengan pustaka utama **[whatsmeow](https://github.com/tulir/whatsmeow)** dan penyimpanan sesi/data berbasis **SQLite**. 

Proyek ini merupakan porting modern dari bot WhatsApp berbasis JavaScript/Node.js (Baileys) demi performa maksimal, efisiensi resource, dan tampilan UI pesan interaktif yang kekinian.

---

## 🚀 Mengapa Memilih Go (whatsmeow) dibanding Node.js (Baileys)?

* **Super Hemat RAM**: Bot NodeJS Baileys umumnya memakan **300MB - 600MB RAM**. Onah-Go hanya memerlukan sekitar **20MB - 50MB RAM**! Sangat ideal untuk VPS murah spek terendah sekalipun.
* **Instan & Responsif**: Go dicompile langsung ke bahasa mesin (native binary) sehingga eksekusi command bot jauh lebih cepat dan minim latensi.
* **Tombol Interaktif Modern**: Menggunakan protokol **WhatsApp Interactive Message (Native Flow)** terbaru, sehingga tombol/list menu akan merender secara sempurna di Android, iOS, dan WhatsApp Web.

---

## 🛠️ Persyaratan Sistem

Sebelum menjalankan bot, pastikan sistem Anda telah memiliki:
1. **Golang** (Versi 1.21 ke atas).
2. **Kompiler C (GCC)** karena SQLite (`go-sqlite3`) menggunakan CGO untuk kompilasi:
   * **Di Windows (Lokal)**: Unduh dan letakkan portabel GCC compiler [w64devkit](https://github.com/skeeto/w64devkit/releases) di folder Anda (misal: `C:\w64devkit`).
   * **Di Linux VPS**: Cukup instal paket `build-essential` bawaan Linux:
     ```bash
     sudo apt update && sudo apt install build-essential -y
     ```

---

## ⚙️ Konfigurasi (`config.json`)

Salin konfigurasi awal Anda di file `config.json` di direktori utama:
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
> **⚠️ PENTING**: Jangan pernah mengunggah (push) file `config.json` yang berisi API Key asli ke GitHub publik demi keamanan akun Anda.

---

## 🚀 Cara Menjalankan Bot di Komputer Lokal

### 1. Unduh Dependensi
Jalankan perintah berikut sekali untuk memasang semua library yang dibutuhkan:
```bash
go mod tidy
```

### 2. Jalankan Program (Development)
* **Windows (Menggunakan w64devkit di C:\w64devkit)**:
  ```powershell
  $env:PATH = 'C:\w64devkit\bin;' + $env:PATH; go run main.go
  ```
* **Linux VPS / MacOS**:
  ```bash
  go run main.go
  ```

Saat dijalankan untuk pertama kali, bot akan menanyakan metode login secara interaktif di terminal Anda:
1. Ketik **`1`** untuk **Pairing Code** (Tautkan menggunakan kode nomor HP yang diset di `config.json`).
2. Ketik **`2`** untuk **QR Code** (Scan barcode langsung dari terminal).

---

## 📦 Panduan Build & Deploy 24/7 di VPS Linux

### Langkah 1: Compile Program Ke Single Binary
Kompilasi program ke bentuk executable tunggal agar bisa dipindahkan dan dijalankan tanpa butuh source code lagi:
* **Di Linux VPS**:
  ```bash
  go build -o onah-bot main.go
  ```
* **Di Windows (untuk dijalankan di Windows lokal)**:
  ```powershell
  $env:PATH = 'C:\w64devkit\bin;' + $env:PATH; go build -o onah-bot.exe main.go
  ```

### Langkah 2: Jalankan Menggunakan PM2 (Rekomendasi)
PM2 sangat andal untuk menjaga bot Anda tetap hidup 24/7 di VPS.

```bash
# Jalankan file binary bot di PM2
pm2 start ./onah-bot --name onah-go

# Melihat log bot realtime
pm2 logs onah-go

# Mematikan bot sementara
pm2 stop onah-go

# Menghidupkan ulang bot
pm2 restart onah-go

# Melihat daftar aplikasi PM2 yang berjalan
pm2 status
```

---

## 📂 Panduan Menulis Fitur Baru (Plugin System)

Bot ini menggunakan sistem **Modular Plugin** otomatis. Anda bisa menambah perintah baru dengan membuat file `.go` baru di dalam folder `/plugins`.

### Contoh Template Plugin Baru (`plugins/hello_plugin.go`):
```go
package plugins

import (
	"strings"
)

func init() {
	Register(Command{
		Name:      "halo",                     // Nama command utama (.halo)
		Tags:      []string{"fun"},            // Kategori menu
		Help:      "Menyapa balik pengguna",    // Deskripsi menu
		Limit:     false,                      // Set true jika butuh sisa limit untuk jalan
		Premium:   false,                      // Set true jika hanya untuk member premium
		OwnerOnly: false,                      // Set true jika hanya untuk owner bot
		Execute: func(ctx *Context) error {
			// Mengirim balasan teks sederhana
			return ctx.Reply("Halo juga, " + ctx.PushName + "! Ada yang bisa saya bantu?")
		},
	})
}
```

---

## 📂 Struktur File Utama
* `main.go`: Bootstrapper utama program.
* `config/`: Kode parsing file `config.json`.
* `database/`: Controller data pengguna (Limit & Balance) menggunakan SQLite.
* `whatsapp/`: Logic penanganan socket `whatsmeow` & interceptor pesan masuk.
* `plugins/`: Tempat menulis modul-modul fitur command bot.
