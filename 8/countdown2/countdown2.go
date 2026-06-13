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
	select {
	case <-time.After(10 * time.Second):
		// Ничего не делаем.
	case <-abort:
		fmt.Println("Запуск отменен!")
		return
	}
	launch()
}

func launch() {
	fmt.Println("Launch")
}
