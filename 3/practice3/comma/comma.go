package main

import (
	"fmt"
	"gobook/3/comma"
)

func main() {
	fmt.Println(comma.Comma("1234567890"))

	fmt.Println(comma.CommaBuff("1234567890"))

	fmt.Println(comma.CommaSlice("1234567890"))
}
