package main

import (
	"fmt"
	popcount "gobook/9/practice9/9.2"
)

func main() {
	fmt.Println(popcount.PopCount(255)) // 8
	fmt.Println(popcount.PopCount(0))   // 0
	fmt.Println(popcount.PopCount(123)) // 6
}
