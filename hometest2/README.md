# Backend Home Test #2

Terapkan layanan web service sederhana menggunakan golang, dengan persyaratan berikut:

### 1. Tetapkan fungsi bernama SingleFizzBuzz yang memiliki perilaku berikut:
1. Fungsi tersebut akan menerima bilangan bulat n. Secara default, fungsi tersebut mengembalikan bilangan bulat n tanpa operasi apa pun
2. Jika n habis dibagi 3, kembalikan Fizz
3. Jika n habis dibagi 5, kembalikan Buzz
4. Jika n habis dibagi 3 dan 5, kembalikan FizzBuzz

### 2. Tetapkan enpoint HTTP bernama GET /range-fizzbuzz yang memiliki persyaratan berikut:
1. Memiliki 2 parameter yang disebut `from` dan `to` , keduanya harus berupa bilangan bulat, dan `from` <= `to`
2. Respons harus mengembalikan nilai SingleFizzBuzz untuk setiap bilangan bulat antara from dan to inklusif (`from` dan `to` disertakan), dibatasi oleh spasi.

### 3. Dapat mencatat permintaan, respons, dan keterlambatannya ke STDOUT.

### 4. Ada beberapa persyaratan kinerja titik akhir:
1. 1 detik sebagai batas waktu
2. Dapat membuat maksimal 1000 goroutine untuk perhitungan pada saat yang sama untuk semua permintaan
3. Menerima maksimal 100 angka sebagai rentang

### 5. Dapat diakhiri dengan baik menggunakan SIGTERM.

# Cara penggunaan
1. Atur port dengan flag, default adalah 8080. Misalnya: `-port 1234`
2. Jika ingin menguji pengendali batas waktu, gunakan flag `-simulateTimeout` dengan nilai dalam milidetik