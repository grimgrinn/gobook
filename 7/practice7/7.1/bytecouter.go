package main

import (
	"bufio"
	"bytes"
	"fmt"
)

type WordCounter int

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
}
