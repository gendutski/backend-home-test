package main

import (
	"fmt"
	"math/rand"
)

func CreateRandomSequence(seq []int, used map[int]bool) []int {
	// jika urutan sudah memiliki 10 angka, rekursi berhenti
	if len(seq) == 10 {
		return seq
	}

	// hasilkan angka acak antara 1 dan 10
	num := rand.Intn(10) + 1

	// jika angka belum digunakan, tambahkan ke urutan
	if !used[num] {
		seq = append(seq, num)
		used[num] = true
	}

	// lanjutkan rekursi untuk menambahkan angka ke urutan
	return CreateRandomSequence(seq, used)
}

func main() {
	used := make(map[int]bool)
	seq := []int{}

	seq = CreateRandomSequence(seq, used)
	fmt.Println(seq)
}
