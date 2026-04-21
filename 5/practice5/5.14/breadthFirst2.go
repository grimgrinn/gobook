package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func processItem(path string) []string {
	var result []string
	things, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error open dir %v\n", err)
		return nil
	}
	for _, thing := range things {
		fullPath := filepath.Join(path, thing.Name())
		fmt.Println(fullPath)
		if thing.IsDir() {
			result = append(result, fullPath)
		}
	}
	return result
}

func main() {
	root := "."

	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	breadthFirst(processItem, []string{root})
}
