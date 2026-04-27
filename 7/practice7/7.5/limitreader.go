package main

import (
	"fmt"
	"io"
	"strings"
)

type limitReader struct {
	r io.Reader
	n int64
}

func (lr *limitReader) Read(p []byte) (int, error) {
	if lr.n <= 0 {
		return 0, io.EOF
	}

	if len(p) > int(lr.n) {
		p = p[:lr.n]
	}

	n, err := lr.r.Read(p)
	lr.n -= int64(n)
	return n, err
}

func LimitReader(r io.Reader, n int64) io.Reader {
	return &limitReader{r: r, n: n}
}

func main() {
	r := strings.NewReader("abcdefghij")
	lr := LimitReader(r, 5)

	buf := make([]byte, 3)
	n, err := lr.Read(buf)
	fmt.Printf("%q, %d, %v\n", buf[:n], n, err) // "abc", 3, nil

	n, err = lr.Read(buf)
	fmt.Printf("%q, %d, %v\n", buf[:n], n, err) // "de", 2, nil

	n, err = lr.Read(buf)
	fmt.Printf("%d, %v\n", n, err) // 0, EOF
}
