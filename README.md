# Bot Telegram Peringatan Dini Cuaca BMKG

Aplikasi ini adalah bot Telegram berbasis Go yang berfungsi untuk memantau dan mendistribusikan informasi peringatan dini cuaca secara waktu nyata (real-time) yang bersumber dari Badan Meteorologi, Klimatologi, dan Geofisika (BMKG).

## Atribusi Sumber Data

Seluruh data cuaca dan peringatan dini yang didistribusikan oleh aplikasi ini bersumber langsung dari portal Data Terbuka BMKG. Aplikasi ini dikembangkan dengan mematuhi pedoman penggunaan data yang ditetapkan oleh penyedia data.

- **Sumber Data Utama**: [Badan Meteorologi, Klimatologi, dan Geofisika (BMKG)](https://data.bmkg.go.id/)
- **Referensi Format Data**: [Repository infoBMKG/data-cap](https://github.com/infoBMKG/data-cap)

## Fitur Utama

1.  **Pemantauan Real-time**: Aplikasi melakukan pemantauan berkala terhadap umpan data (feed) peringatan dini BMKG untuk mendeteksi pembaruan terkini.
2.  **Notifikasi Visual**: Menyertakan infografis resmi dari BMKG dalam pesan notifikasi apabila tersedia, memberikan informasi visual yang lebih jelas kepada pengguna.
3.  **Penyaring Wilayah (Filter)**: Mendukung konfigurasi untuk memfilter notifikasi berdasarkan provinsi tertentu, sehingga pengguna hanya menerima informasi yang relevan.
4.  **Mekanisme Percobaan Ulang (Auto-Retry)**: Dilengkapi dengan mekanisme penanganan kesalahan jaringan yang secara otomatis mencoba kembali permintaan data apabila terjadi kegagalan koneksi.
5.  **Efisiensi Penyimpanan**: Menggunakan sistem penyimpanan berbasis berkas (file-based persistence) yang ringan untuk mencegah duplikasi notifikasi.

## Prasyarat Sistem

- **Go**: Versi 1.21 atau yang lebih baru.
- **Koneksi Internet**: Diperlukan untuk mengakses API BMKG dan API Telegram.

## Konfigurasi

Sebelum menjalankan aplikasi, konfigurasi lingkungan kerja diperlukan melalui berkas `.env`. Buatlah berkas `.env` dengan parameter berikut:

```env
TELEGRAM_TOKEN=token_bot_telegram_anda
TARGET_CHAT_ID=id_chat_tujuan_notifikasi
FILTER_PROVINCE=Jawa Barat, DKI Jakarta, Bali
```

### Keterangan Parameter:
- `TELEGRAM_TOKEN`: Token akses bot yang didapatkan dari BotFather.
- `TARGET_CHAT_ID`: ID numerik (Chat ID) pengguna atau grup Telegram tujuan pengiriman notifikasi.
- `FILTER_PROVINCE` (Opsional): Daftar nama provinsi yang ingin dipantau, dipisahkan dengan koma. Kosongkan jika ingin menerima notifikasi dari seluruh Indonesia.

## Cara Penggunaan

1.  Pastikan dependensi telah terunduh:
    ```bash
    go mod tidy
    ```

2.  Jalankan aplikasi:
    ```bash
    go run cmd/main.go
    ```

3.  Aplikasi akan secara otomatis membuat berkas `history.json` untuk menyimpan riwayat peringatan yang telah diproses.

## Struktur Direktori

- `cmd/`: Berisi titik masuk (entry point) aplikasi.
- `internal/`: Berisi logika inti aplikasi yang terbagi menjadi modul-modul:
    - `bmkg`: Klien HTTP untuk berinteraksi dengan layanan data BMKG.
    - `bot`: Layanan penghubung dengan API Telegram.
    - `model`: Definisi struktur data XML (CAP).
    - `storage`: Manajemen penyimpanan riwayat data.

## Lisensi dan Penafian

Aplikasi ini disediakan sebagai alat bantu distribusi informasi. Pengguna wajib tetap merujuk pada kanal informasi resmi BMKG untuk keputusan-keputusan krusial terkait keselamatan. Pengembang aplikasi tidak bertanggung jawab atas keterlambatan atau ketidakakuratan data yang mungkin terjadi akibat gangguan jaringan atau perubahan format data dari sumber aslinya.
