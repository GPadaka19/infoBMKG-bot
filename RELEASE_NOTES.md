# Release Notes

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

### Deployment

```bash
docker compose up -d
```
