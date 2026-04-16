package main

import (
	"fmt"
	"strings"
	"unicode"
)

func expand(s string, f func(string) string) string {
	var buf strings.Builder
	i := 0
	n := len(s)

	for i < n {
		if s[i] == '$' && i+1 < n && unicode.IsLetter(rune(s[i+1])) {
			i++ // пропускаем $
			start := i
			for i < n && unicode.IsLetter(rune(s[i])) {
				i++
			}
			name := s[start:i]
			buf.WriteString(f(name))
		} else {
			buf.WriteByte(s[i])
			i++
		}
	}
	return buf.String()
}

func main() {
	result := expand("Hello, $name! Today is $day", func(s string) string {
		m := map[string]string{"name": "Jopa", "day": "Jopnik"}
		return m[s]
	})
	fmt.Println(result)
}
