package main

import (
	"fmt"
	"gobook/2/popcount"
)

func main() {
	a := uint64(50)
	fmt.Println(popcount.PopCount(a))
}
