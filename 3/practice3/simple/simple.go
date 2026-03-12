package main

import (
	"fmt"
)

func main() {
	s := "string"
	fmt.Println("строка s - "+s+"\n длина строки s -", len(s))
	fmt.Println("первый символ строки  s ", s[0])

	for i, r := range s {
		fmt.Printf("индекс: %d, руна: %c \n", i, r)
	}
}
