package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <url1> [<url2>...]\n", os.Args[0])
		os.Exit(1)
	}

	urls := os.Args[1:]
	result := fetchFirst(urls)
	fmt.Println("Получен ответ от:", result)
}

func fetchFirst(urls []string) string {
	cancel := make(chan struct{})
	response := make(chan string)

	for _, url := range urls {
		go func(u string) {
			body, err := fetch(u, cancel)
			if err != nil {
				return
			}
			select {
			case response <- u + "; " + body:
			case <-cancel:
			}
		}(url)
	}

	// Ждем первый ответ
	first := <-response
	// Отменяем все остальные запросы
	close(cancel)
	return first
}

func fetch(url string, cancel <-chan struct{}) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Cancel = cancel

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
