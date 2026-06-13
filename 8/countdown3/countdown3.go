package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1)) // Чтение одного байта
		abort <- struct{}{}
	}()

	fmt.Println("Начинаю отсчет. Нажмите <Enter> для отказа.")
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		select {
		case <-tick:
			// Ничего не делаем.
		case <-abort:
			fmt.Println("Запуск отменен!")
			return
		}
	}
	launch()
}

func launch() {
	fmt.Println("Запуск!")
}
