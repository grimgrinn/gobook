package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "a", "link":
			if href := getAttr(n, "href"); href != "" {
				links = append(links, href)
			}
		case "img", "script":
			if src := getAttr(n, "src"); src != "" {
				links = append(links, src)
			}
		}
	}
	// for c := n.FirstChild; c != nil; c = c.NextSibling {
	// 	links = visit(links, c)
	// }

	if n.FirstChild != nil {
		links = visit(links, n.FirstChild)
	}
	if n.NextSibling != nil {
		links = visit(links, n.NextSibling)
	}
	return links
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "visit2: parsing error %v", err)
		os.Exit(1)
	}

	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}
