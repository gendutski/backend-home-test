# Backend Home Test #4

Tujuannya adalah untuk membangun backend untuk dompet elektronik sederhana

Dompet akan diimplementasikan dalam bahasa golang.

### Berikut ini adalah operasi utama yang diharapkan dari dompet:

1. Mendaftarkan pengguna baru
2. Membaca saldo
3. Mengisi saldo
4. Transfer uang antar dompet
5. Mencantumkan N transaksi teratas berdasarkan nilai per pengguna
6. Mencantumkan keseluruhan N pengguna yang bertransaksi teratas berdasarkan nilai

### Berikut ini adalah kriteria evaluasinya:

1. Desain kode dan solusi
2. Standar pengkodean, keterbacaan, dan kualitas
3. Ketepatan
4. Kasus uji
5. Menerapkan penyimpanan data dalam memori dan pengoptimalan kinerja, desain tingkat lanjut, dan pola arsitektur.


## Prasyarat

- [Go](https://golang.org/doc/install) 1.16 or later

## Memulai

### 1. Instal Ketergantungan
Pastikan semua modul Go yang diperlukan telah diinstal:

```bash
go mod tidy
```

### 2. Menjalankan Aplikasi
Untuk menjalankan aplikasi, gunakan perintah berikut (port default yang digunakan adalah 8080 sebagai port http):
```bash
go run main.go
```
for use defined port, use flag `-port`
```
go run main.go -port 1234
```

### 3. Membangun Aplikasi
Untuk membangun aplikasi, gunakan perintah berikut:
```bash
go build -o hometest-4
```

### 4. Menjalankan Biner Aplikasi
Setelah membangun aplikasi, Anda dapat menjalankan biner dengan:
```bash
./hometest-4
```
untuk menggunakan port yang ditentukan, gunakan flag `-port`
```bash
./hometest-4 -port 1234
```

### 5. Menjalankan Uji Unit
Uji unit sangat penting untuk memastikan keandalan aplikasi. Untuk menjalankan uji unit, gunakan perintah berikut:
```bash
go test -race ./...
```

### 6. Cakupan Kode
Untuk memeriksa cakupan pengujian aplikasi ini, gunakan perintah berikut:
```bash
go test -cover ./...
```

