package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const defaultFileRead string = "testfile.txt"

func main() {
	var fileName string
	fmt.Printf("Masukkan file untuk dibaca (default %s): ", defaultFileRead)

	// get file path
	scan := bufio.NewScanner(os.Stdin)
	if !scan.Scan() {
		log.Fatal("Failed to scan input!")
	}
	fileName = strings.TrimSpace(scan.Text())
	if fileName == "" {
		fileName = defaultFileRead
	}

	totalNumbers, sumNumbers, err := CountIntegerInFile(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Total angka pada file:", totalNumbers)
	fmt.Println("Jumlah semua angka:", sumNumbers)
}

func CountIntegerInFile(fileName string) (totalNumbers int, sumNumbers int, err error) {
	// open file
	data, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	// get numbers
	numbers := getNumberInString(string(data))
	totalNumbers = len(numbers)
	for _, n := range numbers {
		sumNumbers += n
	}
	return
}

func getNumberInString(data string) []int {
	var stringHolder string
	var result []int
	for _, r := range data {
		if unicode.IsNumber(r) {
			stringHolder += string(r)
		} else {
			if stringHolder == "" {
				continue
			}
			if i, ok := stringToIntger(stringHolder); ok {
				result = append(result, i)
			}
			stringHolder = ""
		}
	}
	// check if stringHolder not empty
	if stringHolder != "" {
		if i, ok := stringToIntger(stringHolder); ok {
			result = append(result, i)
		}
	}

	return result
}

func stringToIntger(data string) (int, bool) {
	result, err := strconv.Atoi(data)
	if err != nil {
		log.Println("Error convert string to integer:", err.Error())
		return 0, false
	}
	return result, true
}
