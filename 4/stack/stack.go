package main

import "fmt"

func main() {
	stack := []int{}

	stack = append(stack, 5) // Внесение 5 в стек

	fmt.Println(stack)

	stack = append(stack, 2, 3, 4, 5, 6, 6)

	top := stack[len(stack)-1] // Верршинa стека

	fmt.Println(top)

	stack = stack[:len(stack)-1] // Удаление элемента из стека

	fmt.Println(stack)

	s := []int{5, 6, 7, 8, 9}

	s1 := remove(s, 2)

	fmt.Println(s1)

	fmt.Println(s)

	fmt.Println(remove(s, 2))

}

func remove(slice []int, i int) []int {
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}
