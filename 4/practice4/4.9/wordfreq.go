package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	freq := make(map[string]int)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		freq[word]++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка при чтении: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("word\tcount")
	for w, n := range freq {
		fmt.Printf("%s\t%d\n", w, n)
	}
}
