package main

import (
	"flag"
	"fmt"
	"gobook/5/links"
	"log"
	"os"
)

type item struct {
	url   string
	depth int
}

var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)
	tokens <- struct{}{}
	list, err := links.Extract(url)
	<-tokens
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	maxDepth := flag.Int("depth", 0, "максимальная глубина (0 = без ограничений)")
	flag.Parse()

	worklist := make(chan []item)
	var n int

	n++
	go func() {
		start := []item{{os.Args[1], 0}}
		worklist <- start
	}()

	seen := make(map[string]bool)

	for ; n > 0; n-- {
		list := <-worklist
		for _, it := range list {
			if !seen[it.url] {
				seen[it.url] = true
				if *maxDepth > 0 && it.depth >= *maxDepth {
					continue
				}
				n++
				go func(it item) {
					newUrls := crawl(it.url)
					var newItems []item
					for _, url := range newUrls {
						newItems = append(newItems, item{url, it.depth + 1})
					}
					worklist <- newItems
				}(it)
			}
		}
	}
}
