package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

type Node interface {
	String() string
}

type CharData string

func (c CharData) String() string {
	return string(c)
}

type Element struct {
	Name     string
	Attrs    []xml.Attr
	Children []Node
}

func (e *Element) String() string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "<%s", e.Name)
	for _, a := range e.Attrs {
		fmt.Fprintf(&buf, " %s=%q", a.Name.Local, a.Value)
	}
	fmt.Fprintf(&buf, ">")
	for _, c := range e.Children {
		fmt.Fprint(&buf, c.String())
	}
	fmt.Fprintf(&buf, "</%s>", e.Name)
	return buf.String()
}

func parse(r io.Reader) (Node, error) {
	dec := xml.NewDecoder(r)
	var stack []*Element
	var root Node

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			// Создаем новый элемент
			el := &Element{
				Name:  tok.Name.Local,
				Attrs: tok.Attr,
			}
			if len(stack) == 0 {
				root = el
			} else {
				// Добавляем текущий элемент как дочерний к верхушке стека
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, el)
			}
			stack = append(stack, el)
		case xml.EndElement:
			// Заканчиваем текущий элемент
			stack = stack[:len(stack)-1]
		case xml.CharData:
			// Добавляем теуст в текущий элемент
			text := strings.TrimSpace(string(tok))
			if text != "" && len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, CharData(text))
			}
		}
	}
	return root, nil
}

func main() {
	root, err := parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(root.String())
}
