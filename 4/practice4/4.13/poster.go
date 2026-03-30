package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type MovieResponse struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Poster   string `json:"Poster"`
	Response string `json:"Response"`
	Error    string `json:"Error"`
}

func main() {
	title := strings.Join(os.Args[1:], " ")
	if title == "" {
		fmt.Fprintf(os.Stderr, "использование: %s <название фильма>\n", os.Args[0])
		os.Exit(1)
	}

	apiKey := os.Getenv("omdb_api")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "необходимо указать API ключ omdb_api")
		os.Exit(1)
	}

	url := fmt.Sprintf("http://www.omdbapi.com/?t=%s&apikey=%s", strings.ReplaceAll(title, " ", "+"), apiKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка запроса: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка чтения: %v\n", err)
		os.Exit(1)
	}

	var movie MovieResponse
	if err := json.Unmarshal(body, &movie); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка парсина JSON: %v\n", err)
		os.Exit(1)
	}

	if movie.Response != "True" {
		fmt.Fprintf(os.Stderr, "Фильм не найден: %s\n", movie.Error)
		os.Exit(1)
	}

	if movie.Poster == "" || movie.Poster == "N/A" {
		fmt.Fprintf(os.Stderr, "Постер не найден для фильма %s\n", movie.Title)
		os.Exit(1)
	}

	fmt.Printf("Фильм: %s (%s)\n", movie.Title, movie.Year)
	fmt.Printf("Постер: %s\n", movie.Poster)

	posterResp, err := http.Get(movie.Poster)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка скачивания постера: %v\n", err)
		os.Exit(1)
	}

	posterDir := "posters"
	if err := os.MkdirAll(posterDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания папки %s: %v\n", posterDir, err)
		os.Exit(1)
	}
	filename := fmt.Sprintf("%s/%s_%s.jpg", posterDir, strings.ReplaceAll(movie.Title, " ", "_"), movie.Year)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания файла: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = io.Copy(file, posterResp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка сохранения: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Постер сохранен в %s\n", filename)
}
