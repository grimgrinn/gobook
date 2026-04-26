package main

import (
	"bytes"
	"fmt"
	"io"
)

type countingWriter struct {
	iow     io.Writer
	counter int64
}

func (s *countingWriter) Write(p []byte) (n int, err error) {
	n, err = s.iow.Write(p)
	s.counter += int64(n)
	return n, nil
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	cw := &countingWriter{iow: w}
	return cw, &cw.counter
}

func main() {
	var buf bytes.Buffer

	w, count := CountingWriter(&buf)
	fmt.Fprintf(w, "hello")
	fmt.Println(*count) // 5
	fmt.Fprintf(w, " world")
	fmt.Println(*count) // 11
}
