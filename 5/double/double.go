package main

import "fmt"

func double(x int) int {
	return x + x
}

func double1(x int) (result int) {
	defer func() { fmt.Printf("double(%d) = %d\n", x, result) }()
	return x + x
}

func main() {
	// fmt.Println(double(4))
	// fmt.Println(double1(4))
	double1(12)
}
