package main

import (
	"fmt"
	basename "gobook/3/basename1"
	"gobook/3/basename2"
)

func main() {
	fmt.Println(basename.Basename("a/b/c/go"))
	fmt.Println(basename.Basename("c.d.go"))
	fmt.Println(basename.Basename("abc"))

	fmt.Println("===============")

	fmt.Println(basename2.Basename("basename2 a/b/c/go"))
	fmt.Println(basename2.Basename("c.d.go"))
	fmt.Println(basename2.Basename("abc"))
}
