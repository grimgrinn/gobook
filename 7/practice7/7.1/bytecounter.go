package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type WordCounter int
type LineCounter int

func (l *LineCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanLines)
	count := 0
	for scanner.Scan() {
		count++
	}
	*l += LineCounter(count)
	return len(p), nil
}

func (w *WordCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}
	*w += WordCounter(count)
	return len(p), nil
}

func main() {
	var w WordCounter

	w.Write([]byte("hello world"))
	w.Write([]byte("foo bar baz"))
	fmt.Println(w) // 5

	var l LineCounter

	l.Write([]byte("line1 \n line2 \n line3 \n line4 "))
	fmt.Println(l)
}
