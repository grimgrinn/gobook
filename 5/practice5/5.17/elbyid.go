package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func contains(tag string, names []string) bool {
	for _, name := range names {
		if tag == name {
			return true
		}
	}
	return false
}

func ElementsByTagName(doc *html.Node, names ...string) []*html.Node {
	var result []*html.Node

	var visit func(n *html.Node)
	visit = func(n *html.Node) {
		if n.Type == html.ElementNode && contains(n.Data, names) {
			result = append(result, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			visit(c)
		}
	}
	visit(doc)
	return result
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error %v\n", err)
		os.Exit(1)
	}
	images := ElementsByTagName(doc, "img")
	fmt.Printf("Found %d images\n", len(images))

	headings := ElementsByTagName(doc, "h1", "h2", "h3", "h4")

	for _, h := range headings {
		fmt.Printf("Heading: %s\n", h.Data)
	}
}
