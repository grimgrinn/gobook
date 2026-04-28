package main

import (
	"flag"
	"fmt"
	"gobook/7/tempconv"
)

var temp = tempconv.CelsiusFlag("temp", 20.0, "температура")

func main() {
	flag.Parse()
	fmt.Println(*temp)
}
