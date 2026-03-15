package main

import (
	"fmt"
	"os"
)

func main() {
	s1 := os.Args[1]
	s2 := os.Args[2]
	result := IsAnagramma(s1, s2)
	fmt.Println(result)
}

func IsAnagramma(s1 string, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}

	r1 := []rune(s1)
	r2 := []rune(s2)

	r1 = sort(r1)
	r2 = sort(r2)

	for i := range r1 {
		if r1[i] != r2[i] {
			return false
		}
	}

	return true
}

func sort(r []rune) []rune {
	for i := 0; i < len(r); i++ {
		for j := len(r) - 1; j > i; j-- {
			if r[j] < r[j-1] {
				temp := r[j]
				r[j] = r[j-1]
				r[j-1] = temp
			}
		}
	}

	return r
}
