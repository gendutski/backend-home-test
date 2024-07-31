package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var input string
	fmt.Print("Masukkan waktu dalam format 12 jam (HH:MM:SS AM/PM): ")
	scan := bufio.NewScanner(os.Stdin)

	if !scan.Scan() {
		log.Fatal("Failed to scan input!")
	}
	input = scan.Text()

	result, err := ConvertTo24HoursTime(input)
	if err != nil {
		fmt.Println("Invalid input format")
		return
	}
	fmt.Println("Waktu dalam format 24 jam:", result)
}

func ConvertTo24HoursTime(input string) (string, error) {
	t, err := time.Parse("03:04:05 PM", input)
	if err != nil {
		return "", err
	}
	return t.Format("15:04:05"), nil
}
