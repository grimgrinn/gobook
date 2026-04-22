package main

import (
	"fmt"
	"os"
)

func max(args ...int) int {
	max := args[0]

	if len(args) == 0 {
		fmt.Println("enter agrument")
		os.Exit(1)
	}

	for _, n := range args {
		if n > max {
			max = n
		}
	}
	return max
}

func min(args ...int) int {
	min := args[0]

	if len(args) == 0 {
		fmt.Println("enter agrument")
		os.Exit(1)
	}

	for _, n := range args {
		if n < min {
			min = n
		}
	}
	return min
}

func main() {

	fmt.Println(max(1, 6, -2, 15, 9, 200))
	fmt.Println(min(1, 6, -2, 15, 9, 200))
	fmt.Println(max(-5, -10, -3)) // -3
	fmt.Println(min(-5, -10, -3)) // -10
}
