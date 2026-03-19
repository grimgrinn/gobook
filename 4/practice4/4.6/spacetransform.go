package main

import (
	"fmt"
	"unicode"
)

func spaceTransform(s string) string {

	if len(s) <= 1 {
		return s
	}

	write := 0

	str := []rune(s)
	space := false

	for read := 0; read < len(str); read++ {
		if unicode.IsSpace(str[read]) {
			if !space {
				str[write] = ' '
				write++
				space = true
			}
			// если уже пробел пропускаем
		} else {
			str[write] = str[read]
			write++
			space = false
		}
	}
	return string(str[:write])
}

func main() {
	str := "      what a fuck pepe    shneyne  pepe wtfaa   wfaa pepe sneyne"

	fmt.Println(str)

	str = spaceTransform(str)

	fmt.Println(str)

}
