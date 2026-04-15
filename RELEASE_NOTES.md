# Release Notes

<<<<<<< Updated upstream
## v1.1.0 — Link-Based Deduplication Fix (2026-02-10)

### Bug Fix

- **Fixed duplicate alert notifications**: BMKG GUID mengandung timestamp jam yang berubah setiap kali feed di-refresh, menyebabkan alert yang sama dikirim berulang kali. Unique key sekarang menggunakan alert **Link** (contoh: `CYG20260209009_alert.xml`) yang stabil per event.

### Breaking Changes

- Format `history.json` berubah dari GUID-based ke Link-based. Pada deploy pertama setelah update, alert yang masih aktif di RSS feed akan dikirim ulang sekali.

---

## v1.0.0 — WhatsApp Edition (2026-02-03)

### Highlights

Rilis pertama bot peringatan dini cuaca BMKG dengan distribusi notifikasi melalui **WhatsApp Gateway**.

### Features

- **Real-time Monitoring**: Polling otomatis RSS feed BMKG setiap 3 menit.
- **WhatsApp Notification**: Pengiriman alert ke nomor individu atau grup WhatsApp via WhatsApp Gateway API.
- **Visual Alert**: Menyertakan infografis resmi BMKG dari data CAP (Common Alerting Protocol) jika tersedia.
- **Province Filtering**: Filter notifikasi berdasarkan provinsi melalui environment variable `FILTER_PROVINCE`.
- **Auto-Retry & Self-Healing**: Retry otomatis untuk HTTP request, serta reconnect otomatis jika sesi WhatsApp Gateway terputus.
- **Dry Run Mode**: Mode console-only jika konfigurasi WhatsApp tidak diset, untuk pengujian lokal.
- **Lightweight Persistence**: Penyimpanan riwayat alert berbasis file (`history.json`) tanpa database.
- **Docker Integration**: Multi-stage Dockerfile dan Docker Compose untuk deployment di VPS.

### Limitations

- **Single Broadcast Target**: Saat ini hanya mendukung pengiriman ke satu tujuan (nomor/grup) yang dikonfigurasi di `.env`.
- **No Interactive Commands**: Belum mendukung manajemen multi-user atau perintah interaktif via chat.
=======
## v2.0.0 - WhatsApp Migration (2026-02-07)

### 🚀 Breaking Changes

- **Telegram → WhatsApp**: Bot sekarang mengirim notifikasi via WhatsApp menggunakan gateway ([go-whatsapp-web-multidevice](https://github.com/aldinokemal/go-whatsapp-web-multidevice))
- Environment variables baru (lihat `.env.example`):
  - `WA_GATEWAY_URL` - Endpoint gateway WhatsApp (default: `http://gowa:3006/whatsapp/send/message`)
  - `WA_TARGET_PHONE` - Nomor/Group ID target
  - `WA_USER` / `WA_PASSWORD` - Basic Auth (opsional)
- Environment variables dihapus:
  - `TELEGRAM_TOKEN`
  - `TARGET_CHAT_ID`

### ✨ Features

- Support kirim ke **individual** (`6281xxx`, `6281xxx@s.whatsapp.net`) atau **group** (`120363xxx@g.us`)
- Sanitasi nomor otomatis: `08xx` → `628xx`, `8xx` → `628xx`
- Kirim foto infografik BMKG via `image_url` — gateway download otomatis
- Auto-healing: reconnect otomatis jika sesi WhatsApp terputus (`AUTHENTICATION_ERROR`)
- Fallback ke text jika pengiriman foto gagal
- Caption auto-truncate jika melebihi 1000 karakter

### 🔧 Technical Changes
>>>>>>> Stashed changes

| Component            | Change                                           |
| -------------------- | ------------------------------------------------ |
| `internal/whatsapp/` | ✨ New - WhatsApp gateway client (HTTP + JSON)   |
| `internal/bot/`      | ❌ Deleted - Telegram implementation             |
| `cmd/main.go`        | 🔄 Updated - WhatsApp integration, polling loop  |
| `docker-compose.yml` | 🔄 Updated - External network `infobmkg-bot-net` |
| `Dockerfile`         | 🔄 Updated - Go 1.24, run as root for volume fix |
| `go.mod`             | 🔄 Updated - Removed `telegram-bot-api`, Go 1.24 |

<<<<<<< Updated upstream
```bash
docker compose up -d
```
=======
### 📦 Dependencies

```
github.com/go-resty/resty/v2 v2.17.1   # HTTP client (BMKG fetch)
github.com/joho/godotenv v1.5.1         # .env loader
golang.org/x/net v0.43.0               # indirect
```

### 🏗️ Architecture

```
cmd/main.go              → Entry point, polling loop (3 min interval)
internal/bmkg/client.go  → BMKG RSS + CAP XML fetcher (retry 3x, timeout 15s)
internal/model/types.go  → RSS & CAP data structures
internal/storage/history.go → File-based dedup (history.json)
internal/whatsapp/client.go → WhatsApp gateway client (auto-heal, image/text)
```

### ⚠️ Known Limitations

- **BMKG RSS Feed Delay**: RSS feed BMKG (`rss.xml`) di-rebuild secara batch, bukan real-time. Alert yang diterbitkan BMKG bisa delayed 30–90 menit sebelum muncul di RSS feed. Ini adalah limitasi dari sisi BMKG, bukan bot.
- Polling interval 3 menit — bot sudah optimal untuk menangkap alert segera setelah RSS diperbarui.

### 🐳 Deployment

```bash
# Update .env dengan konfigurasi WhatsApp
# Pastikan network external sudah ada (shared dengan gateway gowa)
docker network create infobmkg-bot-net

# Deploy
docker-compose up -d --build
```

---

## v1.0.0 - Initial Release

- Telegram Bot untuk notifikasi peringatan cuaca BMKG
- Polling RSS feed BMKG secara berkala
- Filter berdasarkan provinsi (`FILTER_PROVINCE`)
- Kirim foto infografik dengan caption
- History tracking via `history.json` untuk mencegah duplikasi
>>>>>>> Stashed changes
