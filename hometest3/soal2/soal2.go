package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	var data string
	fmt.Print("Masukkan kalimat: ")
	scan := bufio.NewScanner(os.Stdin)

	if !scan.Scan() {
		log.Fatal("Failed to scan input!")
	}
	data = scan.Text()
	if IsPalindrom(data) {
		fmt.Println("Merupakan palindrom")
	} else {
		fmt.Println("Bukan palindrom")
	}
}

func IsPalindrom(data string) bool {
	// lowercase string
	data = strings.ToLower(data)

	// strip non word
	reg := regexp.MustCompile(`\W`)
	wordOnly := reg.ReplaceAll([]byte(data), []byte(""))

	var result bool = true
	length := len(wordOnly)
	for i := 0; i < length/2; i++ {
		if wordOnly[i] != wordOnly[(length-i-1)] {
			result = false
			break
		}
	}
	return result
}
