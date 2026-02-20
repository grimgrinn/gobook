package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Printf("Privet, %s", os.Args[1])
	} else {
		fmt.Printf("Ty zabyl predstavitsia")
	}
}
