package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

func main() {
	// ключ - название категории, значение - счетчик
	categories := make(map[string]int)
	invalid := 0

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "catscount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}

		switch {
		case unicode.IsLetter(r):
			categories["letters"]++
		case unicode.IsDigit(r):
			categories["digit"]++
		case unicode.IsSpace(r):
			categories["space"]++
		case unicode.IsPunct(r):
			categories["punct"]++
		case unicode.IsSymbol(r):
			categories["symbol"]++
		default:
			categories["other"]++
		}
	}

	for cat, count := range categories {
		fmt.Printf("%s: %d\n", cat, count)
	}

	if invalid > 0 {
		fmt.Printf("\n%d неверных символов utf-8\n", invalid)
	}

}
