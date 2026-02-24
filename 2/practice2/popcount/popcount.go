package main

import (
	"fmt"
	"gobook/2/popcount"
)

func main() {
	a := uint64(5)
	fmt.Println(popcount.PopCount(a))
	fmt.Println(popcount.PopCountLoop(a))
	fmt.Println(popcount.PopCountBytes(a))
	fmt.Println(popcount.PopCountClear(a))
}
