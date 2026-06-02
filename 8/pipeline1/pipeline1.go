package main

import "fmt"

func main() {

	naturals := make(chan int)
	squares := make(chan int)

	// Генерация
	go func() {
		for x := 0; ; x++ {
			naturals <- x
		}
	}()

	// Возведение в квадрат
	go func() {
		for {
			x := <-naturals
			squares <- x * x
		}
	}()

	// Вывод (в главной go-подпрограмме)
	for {
		fmt.Println(<-squares)
	}
}
