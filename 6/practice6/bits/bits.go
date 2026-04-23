package main

import "fmt"

func main() {

	var n int
	// var w uint64 = 0

	// w |= 1 << 2
	// fmt.Printf("%012b\n", w) // 000000000100

	// w |= 1 << 5
	// fmt.Printf("%012b\n", w) // 000000100100

	// fmt.Println((w & (1 << 2)) != 0) // true
	// fmt.Println((w & (1 << 3)) != 0) // false

	n = 80
	s := make([]uint64, 9)
	fmt.Printf("%d\n%064b\n", n, n)

	fmt.Println(s)
}
