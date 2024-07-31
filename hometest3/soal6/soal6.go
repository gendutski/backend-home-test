package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	fmt.Println(`1. Buatlah program yang mengeluarkan angka dari 1 hingga 100, dengan kondisi jika angka tersebut habis dibagi 3, output adalah "Fizz", jika habis dibagi 5, output adalah "Buzz", jika habis dibagi 3 dan 5, output adalah "FizzBuzz".`)
	fmt.Println(`2. Buatlah program untuk menentukan apakah suatu string adalah palindrom atau tidak.`)
	fmt.Println(`3. Buatlah program untuk menghitung bilangan faktorial dari suatu bilangan.`)
	fmt.Println(`4. Buatlah program untuk mencari nilai terbesar dan terkecil dalam sebuah array integer.`)
	fmt.Println(`5. Buatlah program untuk menghitung jumlah karakter huruf dan angka pada sebuah string.`)
	fmt.Print("Pilih 1-5: ")

	scan := bufio.NewScanner(os.Stdin)
	if !scan.Scan() {
		log.Fatal("Failed to scan input!")
	}
	no, err := strconv.Atoi(scan.Text())
	if err != nil {
		log.Fatal("Invalid input:", err)
	}

	switch no {
	case 1:
		CaseNo1()
	case 2:
		CaseNo2()
	case 3:
		CaseNo3()
	case 4:
		CaseNo4()
	case 5:
		CaseNo5()
	}
}

func CaseNo1() {
	for i := 1; i <= 100; i++ {
		if i%3 == 0 && i%5 == 0 {
			fmt.Print("FizzBuzz ")
		} else if i%3 == 0 {
			fmt.Print("Fizz ")
		} else if i%5 == 0 {
			fmt.Print("Buzz ")
		} else {
			fmt.Print(i, " ")
		}
	}
	fmt.Println()
}

func CaseNo2() {
	var data string
	fmt.Print("Masukkan kalimat: ")
	scan := bufio.NewScanner(os.Stdin)

	if !scan.Scan() {
		log.Fatal("Failed to scan input!")
	}
	data = scan.Text()

	// check palindrom
	// lowercase string
	data = strings.ToLower(data)

	// strip non word
	reg := regexp.MustCompile(`\W`)
	wordOnly := reg.ReplaceAll([]byte(data), []byte(""))

	var isPalindrom bool = true
	length := len(wordOnly)
	for i := 0; i < length/2; i++ {
		if wordOnly[i] != wordOnly[(length-i-1)] {
			isPalindrom = false
			break
		}
	}

	if isPalindrom {
		fmt.Println("Merupakan palindrom")
	} else {
		fmt.Println("Bukan palindrom")
	}
}

func CaseNo3() {
	fmt.Print("Masukkan nilai bilangan untuk medapatkan faktorial-nya: ")

	scan := bufio.NewScanner(os.Stdin)
	if !scan.Scan() {
		log.Fatal("Failed to scan input!")
	}
	val, err := strconv.Atoi(scan.Text())
	if err != nil {
		log.Fatal("Invalid input:", err)
	}

	result := 1
	for i := 2; i <= val; i++ {
		result *= i
	}
	fmt.Printf("Nilai %d! adalah: %d\n", val, result)
}

func CaseNo4() {
	fmt.Print("Masukkan nilai bilangan dipisahkan koma untuk mencari nilai min dan max: ")
	scan := bufio.NewScanner(os.Stdin)
	if !scan.Scan() {
		log.Fatal("Failed to scan input!")
	}

	// get numbers
	var source []int
	for _, s := range strings.Split(scan.Text(), ",") {
		i, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			continue
		}
		source = append(source, i)
	}

	if len(source) < 2 {
		log.Fatalln("Nilai yang dimasukkan harus lebih dari 2 nilai")
	}

	min := source[0]
	max := source[1]
	for i := 1; i < len(source); i++ {
		if source[i] < min {
			min = source[i]
		}
		if source[i] > max {
			max = source[i]
		}
	}
	fmt.Printf("Nilai min: %d, nilai max: %d\n", min, max)
}

func CaseNo5() {
	fmt.Print("Masukkan kalimat: ")
	scan := bufio.NewScanner(os.Stdin)
	if !scan.Scan() {
		log.Fatal("Failed to scan input!")
	}

	var totalNumber, totalCharacter int
	for _, r := range scan.Text() {
		if unicode.IsLetter(r) {
			totalCharacter++
		} else if unicode.IsNumber(r) {
			totalNumber++
		}
	}
	fmt.Printf("Jumlah karakter huruf: %d, jumlah karakter angka: %d\n", totalCharacter, totalNumber)
}
