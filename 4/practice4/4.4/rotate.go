package main

import "fmt"

func rotate(slice []int) {
	if len(slice) <= 1 {
		return
	}
	last := slice[len(slice)-1]
	//first := slice[0]             	  // Влево
	//copy(slice, slice[1:])			  // Влево
	copy(slice[1:], slice[:len(slice)-1]) // Вправо
	slice[0] = last                       // Вправо
}

func main() {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	fmt.Println(slice)
	rotate(slice)
	fmt.Println((slice))
	rotate(slice)
	fmt.Println((slice))
	rotate(slice)
	fmt.Println((slice))
	rotate(slice)
	fmt.Println((slice))
	rotate(slice)
	fmt.Println((slice))
	rotate(slice)
	fmt.Println((slice))
}
