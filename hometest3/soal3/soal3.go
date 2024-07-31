package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const sourceURL string = "https://jsonplaceholder.typicode.com/posts"

type Source struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	source, err := GetData(sourceURL)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := json.MarshalIndent(source, "", "\t")
	if err != nil {
		log.Fatalln("Failed to MarshalIndent:", err)
	}
	fmt.Println(string(result))
}

func GetData(url string) ([]*Source, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []*Source
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
