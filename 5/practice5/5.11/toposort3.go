package main

import (
	"fmt"
	"os"
)

// prereqs отображает каждый курс на список курсов, которые
// должны быть прочитаны раньше него.
// Добавлен цикл по условию
var prereqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},
	"compilers": {"data structures",
		"formal languages",
		"computer organization",
	},
	"data structures": {"discrete math"},
	"databases":       {"data structures"},
	"discrete math":   {"intro to programming", "algorithms"}, // цикл
	//	"discrete math":    {"intro to programming"}, // не цикл
	"formal languages": {"discrete math"},
	"networks":         {"operating systems"},
	"operating systems": {"data structures",
		"computer organization"},
	"programming languages": {"data structures",
		"computer organization"},
}

const (
	visiting = 1
	visited  = 2
)

func topoSort(m map[string][]string) ([]string, error) {
	var order []string

	state := make(map[string]int)
	var visitAll func(items []string) error
	visitAll = func(items []string) error {
		for _, item := range items {
			if state[item] == visited {
				continue
			}
			if state[item] == visiting {
				return fmt.Errorf("цикл: %s", item)
			}
			state[item] = visiting
			if err := visitAll(m[item]); err != nil {
				return err
			}
			state[item] = visited
			order = append(order, item)
		}
		return nil
	}

	for key := range m {
		if err := visitAll([]string{key}); err != nil {
			return nil, err
		}
	}
	return order, nil
}

func main() {
	order, err := topoSort(prereqs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
		os.Exit(1)
	}
	for i, course := range order {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}
