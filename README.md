# BMKG Weather Alert Telegram Bot

Bot Telegram sederhana yang ditulis dengan Go untuk memantau peringatan dini cuaca dari BMKG (Badan Meteorologi, Klimatologi, dan Geofisika) dan mengirimkan notifikasi real-time ke Telegram.

## 🚀 Fitur

- 📡 **Real-time Polling**: Memantau feed RSS BMKG setiap 3 menit.
- ⚡ **Cepat & Ringan**: Ditulis dengan Go, menggunakan Resty Client.
- 💾 **Persistence Sederhana**: Menggunakan file JSON untuk menyimpan riwayat alert (tanpa database berat).
- 🛡️ **Auto-Retry**: Otomatis mencoba ulang jika koneksi ke BMKG gagal.
- 🐳 **Docker Ready**: Siap dijalankan di container (TODO).

## 🛠️ Cara Menjalankan

### Prasyarat
- Go 1.21+

### 1. Set Environment Variables

Anda perlu Token Bot Telegram dan Chat ID tujuan notifikasi.

**Windows (PowerShell):**
```powershell
$env:TELEGRAM_TOKEN = "YOUR_BOT_TOKEN_HERE"
$env:TARGET_CHAT_ID = "YOUR_CHAT_ID_HERE"
go run cmd/main.go
```

**Linux/Mac:**
```bash
export TELEGRAM_TOKEN="YOUR_BOT_TOKEN_HERE"
export TARGET_CHAT_ID="YOUR_CHAT_ID_HERE"
go run cmd/main.go
```

### 2. Mode Dry Run (Tanpa Token)
Jika dijalankan tanpa token, bot akan berjalan dalam mode **Dry Run**, di mana notifikasi hanya akan dicetak ke terminal (console output) untuk tujuan debugging.

```bash
go run cmd/main.go
```

## 📂 Struktur Proyek

```
bmkg-telegram-bot/
├── cmd/
│   └── main.go              # Entry point aplikasi
├── internal/
│   ├── bmkg/                # Client API BMKG (Resty)
│   ├── bot/                 # Wrapper API Telegram
│   ├── model/               # Struktur Data XML
│   └── storage/             # Logika penyimpanan (JSON)
└── history.json             # File penyimpanan riwayat (auto-generated)
```

## ⚠️ Lisensi & Atribut
Data cuaca disediakan oleh [BMKG](https://data.bmkg.go.id/). Bot ini mematuhi aturan atribusi dengan mencantumkan sumber data pada setiap notifikasi.
