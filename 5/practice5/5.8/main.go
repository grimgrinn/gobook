package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func forEachNode(n *html.Node, pre, post func(n *html.Node) bool) bool {
	if pre != nil {
		if !pre(n) {
			return false
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if !forEachNode(c, pre, post) {
			return false
		}
	}

	if post != nil {
		if !post(n) {
			return false
		}
	}

	return true
}

func ElementByID(doc *html.Node, id string) *html.Node {
	var result *html.Node

	// Функция pre: проверяет каждый узел
	pre := func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			fmt.Println("Proveryayu:", n.Data) // otladka
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == id {
					result = n
					return false // останавливаем обход
				}
			}
		}
		return true
	}

	forEachNode(doc, pre, nil)
	return result
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	id := os.Args[1] // передаем искомый id через аргумент
	node := ElementByID(doc, id)

	if node != nil {
		fmt.Printf("Найден элемент: <%s>\n", node.Data)
		// выводим атрибуты
		for _, a := range node.Attr {
			fmt.Printf(" %s=\"%s\"\n", a.Key, a.Val)
		}
	} else {
		fmt.Printf("Элемент с id=%s не найден\n", id)
	}
}
