package main

import (
	"fmt"
	"time"
)

func stage(in <-chan int, out chan<- int, name string) {
	fmt.Println(name)
	for v := range in {
		out <- v * 2
	}
	close(out)
}

func main() {
	const N = 10 // количество этапов

	start := time.Now()

	// Создаем перый канал
	first := make(chan int, 1)
	var prev <-chan int = first

	for i := 0; i < N; i++ {
		next := make(chan int, 1)
		go stage(prev, next, fmt.Sprintf("stage%d", i))
		prev = next
	}

	first <- 1

	result := <-prev
	close(first)

	fmt.Printf("Результат: %d\n", result) // 1 * 2 ^ N
	fmt.Printf("Время транзита: %v\n", time.Since(start))
}
