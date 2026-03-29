package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Comic struct {
	Num         int    `json:"num"`
	Title       string `json:"title"`
	Trancscript string `json:"transcript"`
	Img         string `json:"img"`
	Alt         string `json:"alt"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "использование: %s <поисковый запрос>\n", os.Args[0])
		os.Exit(1)
	}
	query := strings.Join(os.Args[1:], " ")

	// Читаем JSON
	data, err := os.ReadFile("xkcd.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения: %v\n", err)
		os.Exit(1)
	}

	var comics []Comic
	if err := json.Unmarshal(data, &comics); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка парсинга: %v\n", err)
	}

	// Поиск
	found := false
	for _, comic := range comics {
		if strings.Contains(strings.ToLower(comic.Title), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(comic.Trancscript), strings.ToLower(query)) {
			fmt.Printf("https://xkcd.com/%d/\n", comic.Num)
			fmt.Printf("Title: %s\n", comic.Title)
			fmt.Printf("Alt: %s\n", comic.Alt)
			fmt.Println()
			found = true
		}
	}

	if !found {
		fmt.Println("Ничего не найдено")
	}
}
