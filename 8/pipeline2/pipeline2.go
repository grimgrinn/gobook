package main

import "fmt"

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	// Генерация
	go func() {
		for x := 0; x < 100; x++ {
			naturals <- x
		}
		close(naturals)
	}()

	// Возведение в квадрат
	go func() {
		for x := range naturals {
			squares <- x * x
		}
		close(squares)
	}()

	// Вывод (в главной go-подпрограмме)
	for x := range squares {
		fmt.Println(x)
	}
}
