# Bot WhatsApp Peringatan Dini Cuaca BMKG

Aplikasi bot berbasis Go yang memantau dan mendistribusikan informasi peringatan dini cuaca secara real-time dari BMKG melalui WhatsApp.

## Atribusi Sumber Data

Seluruh data cuaca dan peringatan dini yang didistribusikan oleh aplikasi ini bersumber langsung dari portal Data Terbuka BMKG. Aplikasi ini dikembangkan dengan mematuhi pedoman penggunaan data yang ditetapkan oleh penyedia data.

- **Sumber Data Utama**: [Badan Meteorologi, Klimatologi, dan Geofisika (BMKG)](https://data.bmkg.go.id/)
- **Referensi Format Data**: [Repository infoBMKG/data-cap](https://github.com/infoBMKG/data-cap)

## Fitur Utama

1.  **Pemantauan Real-time**: Polling otomatis terhadap RSS feed peringatan dini BMKG setiap 3 menit.
2.  **Notifikasi WhatsApp**: Mengirim peringatan cuaca langsung ke nomor individu atau grup WhatsApp melalui WhatsApp Gateway.
3.  **Notifikasi Visual**: Menyertakan infografis resmi dari BMKG (data CAP) dalam pesan notifikasi apabila tersedia.
4.  **Penyaring Wilayah**: Mendukung filter notifikasi berdasarkan provinsi melalui environment variable.
5.  **Deduplikasi Akurat**: Menggunakan alert link sebagai unique key untuk mencegah duplikasi, karena GUID BMKG mengandung timestamp yang berubah setiap jam.
6.  **Auto-Retry & Self-Healing**: Mekanisme retry otomatis untuk koneksi HTTP, serta reconnect otomatis ke WhatsApp Gateway jika sesi terputus.
7.  **Dry Run Mode**: Jika konfigurasi WhatsApp tidak diset, bot berjalan dalam mode console-only untuk pengujian.
8.  **Docker Ready**: Dikontainerisasi dengan multi-stage Dockerfile dan Docker Compose untuk deployment yang mudah.

## Prasyarat Sistem

- **Go**: Versi 1.24 atau lebih baru.
- **WhatsApp Gateway**: Instansi [go-whatsapp-web-multidevice](https://github.com/AdiKhoworker/go-whatsapp-web-multidevice) yang sudah berjalan.
- **Docker** (opsional): Untuk deployment via container.

## Konfigurasi

Salin `.env.example` menjadi `.env` dan sesuaikan parameternya:

```env
# WhatsApp Gateway Configuration
WA_GATEWAY_URL=http://gowa:3006/whatsapp/send/message
# Untuk individu: 6281xxx atau 6281xxx@s.whatsapp.net
# Untuk grup: 120363xxx@g.us
WA_TARGET_PHONE=6281xxx
WA_USER=
WA_PASSWORD=

# Filter Provinsi (opsional, pisahkan dengan koma)
FILTER_PROVINCE=DI Yogyakarta, Jawa Tengah
```

### Keterangan Parameter

| Parameter | Wajib | Keterangan |
|---|---|---|
| `WA_GATEWAY_URL` | Ya | URL endpoint WhatsApp Gateway (`/send/message`). |
| `WA_TARGET_PHONE` | Ya | Nomor tujuan atau Group ID WhatsApp. |
| `WA_USER` | Tidak | Username Basic Auth untuk gateway (jika ada). |
| `WA_PASSWORD` | Tidak | Password Basic Auth untuk gateway (jika ada). |
| `FILTER_PROVINCE` | Tidak | Daftar provinsi yang dipantau, dipisahkan koma. Kosongkan untuk seluruh Indonesia. |

## Cara Penggunaan

### Lokal

```bash
go mod tidy
go run cmd/main.go
```

### Docker Compose

```bash
docker compose up -d
```

Pastikan WhatsApp Gateway dan bot berada di network Docker yang sama (`infobmkg-bot-net`).

## Struktur Direktori

```
data-cap/
├── cmd/
│   └── main.go              # Entry point, polling loop, & notification logic
├── internal/
│   ├── bmkg/
│   │   └── client.go         # HTTP client untuk RSS feed & CAP XML BMKG
│   ├── model/
│   │   └── types.go          # Definisi struct XML (RSS, CAP Alert)
│   ├── storage/
│   │   └── history.go        # Penyimpanan riwayat alert (file-based)
│   └── whatsapp/
│       └── client.go         # Client WhatsApp Gateway (text & image)
├── Dockerfile                # Multi-stage build
├── docker-compose.yml        # Konfigurasi deployment
├── .env.example              # Template environment variable
└── history.json              # Riwayat alert (auto-generated, gitignored)
```

## Lisensi dan Penafian

Aplikasi ini disediakan sebagai alat bantu distribusi informasi. Pengguna wajib tetap merujuk pada kanal informasi resmi BMKG untuk keputusan krusial terkait keselamatan. Pengembang tidak bertanggung jawab atas keterlambatan atau ketidakakuratan data yang mungkin terjadi akibat gangguan jaringan atau perubahan format data dari sumber aslinya.
