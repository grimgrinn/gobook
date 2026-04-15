package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func countElements(n *html.Node, counts map[string]int) {
	if n.Type == html.ElementNode {
		counts[n.Data]++
	}
	if n.FirstChild != nil {
		countElements(n.FirstChild, counts)
	}
	if n.NextSibling != nil {
		countElements(n.NextSibling, counts)
	}
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошиюка парсинга: %v\n", err)
		os.Exit(1)
	}

	counts := make(map[string]int)
	countElements(doc, counts)

	for tag, cnt := range counts {
		fmt.Printf("%s: %d\n", tag, cnt)
	}
}
