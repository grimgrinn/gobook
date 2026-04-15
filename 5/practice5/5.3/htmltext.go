package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func printText(n *html.Node) {
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
		return
	}
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			fmt.Println(text)
		}
	}
	if n.FirstChild != nil {
		printText(n.FirstChild)
	}
	if n.NextSibling != nil {
		printText(n.NextSibling)
	}
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "htmltext: parsing error %v", err)
		os.Exit(1)
	}

	printText(doc)
}
