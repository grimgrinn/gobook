package main

import (
	"fmt"
	"unicode/utf8"
)

func unicodeReverse(str string) string {
	bytes := []byte(str)

	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	return string(bytes)
}

func reverseUTF8(b []byte) []byte {
	// Сначала переворачиваем руны
	for i := 0; i < len(b); {
		_, size := utf8.DecodeRune(b[i:])
		reverseBytes(b[i : i+size])
		i += size
	}
	// Теперь переворачиваем весь срез
	reverseBytes(b)
	return b
}

func reverseBytes(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func main() {
	s := "это строка символос юникод"

	fmt.Println(s)

	fmt.Println(unicodeReverse(s))

	fmt.Println(string(reverseUTF8([]byte(s))))
}
