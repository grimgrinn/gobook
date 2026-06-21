package main

import (
	"fmt"
	"gobook/5/links"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var tokens = make(chan struct{}, 20)

func crawl(url string, cancel <-chan struct{}) []string {
	fmt.Println(url)

	tokens <- struct{}{}
	defer func() { <-tokens }()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print(err)
		return nil
	}
	req.Cancel = cancel

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return nil
	}
	defer resp.Body.Close()

	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
		return nil
	}
	return list
}

func main() {
	cancel := make(chan struct{})

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		fmt.Println("\nОтмена запросов...")
		close(cancel)
	}()

	worklist := make(chan []string)
	var n int

	n++
	go func() {
		worklist <- os.Args[1:]
	}()

	seen := make(map[string]bool)

	for ; n > 0; n-- {
		select {
		case <-cancel:
			fmt.Println("Программа остановлена")
			return

		case list := <-worklist:
			for _, link := range list {
				if !seen[link] {
					seen[link] = true
					n++
					go func(link string) {
						newLinks := crawl(link, cancel)
						select {
						case <-cancel:
							return
						case worklist <- newLinks:
						}
					}(link)
				}
			}
		}
	}
}
