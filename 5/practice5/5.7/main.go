package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func forEachElement(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachElement(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

var depth int

func startElement(n *html.Node) {

	switch n.Type {
	case html.ElementNode:
		fmt.Printf("%*s<%s", depth*2, "", n.Data)
		for _, a := range n.Attr {
			fmt.Printf(" %s=\"%s\"", a.Key, a.Val)
		}
		if n.FirstChild == nil {
			fmt.Printf("/>\n")
		} else {
			fmt.Printf(">\n")
			depth++
		}
	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" {
			fmt.Printf("%*s%s\n", depth*2, "", text)
		}
	case html.CommentNode:
		fmt.Printf("%*s<!-- %s -->\n", depth*2, "", n.Data)
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode && n.FirstChild != nil {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}

func main() {
	// Тестовый HTML
	testHTML := `
	<html>
	<head><title>Test</title></head>
	<body>
	<!-- comment -->
		<div class="main" id="content">
			<p>Hello</p>
			<img src="pic.jpg"/>
		</div>
	</body>
	</html>
	`

	doc, err := html.Parse(strings.NewReader(testHTML))
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	forEachElement(doc, startElement, endElement)
}
