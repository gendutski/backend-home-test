# Soal 4

Buatlah sebuah program dengan Golang untuk membuat rekursi dan menampilkan angka 1
sampai 10 dalam urutan acak. Program harus menggunakan rekursi untuk menghasilkan
urutan angka yang berbeda pada setiap eksekusi.
Output yang diharapkan (contoh):

```
7 1 5 10 9 8 2 6 3 4
```

My Note:

Cara sederhananya seharusnya seperti dibawah ini:
```
target := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
rand.Shuffle(len(target), func(i, j int) {
    target[i], target[j] = target[j], target[i]
})
fmt.Println(target)
```
Tapi karena harus rekursi jadinya tidak pakai ini