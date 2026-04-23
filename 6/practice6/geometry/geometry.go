package main

import (
	"fmt"
	"gobook/6/geometry"
)

func main() {
	perim := geometry.Path{{1, 1}, {5, 1}, {5, 4}, {1, 1}}
	fmt.Println(geometry.Path.Distance(perim)) // "12", автономная функция
	fmt.Println(perim.Distance())              // "12", метод geometry.Path
}
