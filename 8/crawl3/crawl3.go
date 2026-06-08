package main

import (
	"fmt"
	"gobook/5/links"
	"log"
	"os"
)

func main() {
	worklist := make(chan []string)  // Список URL, могут быть дубли.
	unseenLinks := make(chan string) // Удаление дублей.

	// Добавление в список аргументов командной строки.
	go func() { worklist <- os.Args[1:] }()

	// Создание 20 go-подпрограмм сканирования для выборки
	// всех непросмотренных ссылок.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// Главная go-подпрограмма удаляет дубликаты из списка
	// и отправляет непросмотренные ссылки сканерам.
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}
