package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func countWordsAndImages(n *html.Node) (words, images int) {
	// Пропускаем script и style
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
		return
	}

	// Если текстоывый узел - считаем слова
	if n.Type == html.TextNode {
		words += len(strings.Fields(n.Data))
	}

	if n.Type == html.ElementNode && n.Data == "img" {
		images++
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		w, i := countWordsAndImages(c)
		words += w
		images += i
	}

	return
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	words, images := countWordsAndImages(doc)
	fmt.Printf("Слов: %d, изображений: %d\n", words, images)
}
