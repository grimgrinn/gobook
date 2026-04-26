package main

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

type stringReader struct {
	str      string
	position int
}

func (s *stringReader) Read(p []byte) (n int, err error) {
	if s.position >= len(s.str) {
		return 0, io.EOF
	}

	available := len(s.str) - s.position
	if available > len(p) {
		available = len(p)
	}

	copy(p, s.str[s.position:s.position+available])
	s.position += available
	return available, nil
}

func NewReader(s string) io.Reader {
	return &stringReader{str: s}
}

func main() {
	htmlStr := "<html><body><h1>Hello</h1></body></html>"
	r := NewReader(htmlStr)
	doc, err := html.Parse(r)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}

	fmt.Printf("Parsed: %T\n", doc)
}
