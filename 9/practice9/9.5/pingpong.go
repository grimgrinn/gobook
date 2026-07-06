package main

import (
	"fmt"
	"time"
)

// ping отправляет сообщение в канал и ждет ответа
func ping(pingChan chan string, pongChan chan string, done chan struct{}, n int) {
	for i := 0; i < n; i++ {
		pingChan <- "ping"
		<-pongChan // ждем ответа
	}
	done <- struct{}{}
}

// pong получает сообщение и отправляет обратно
func pong(pingChan chan string, pongChan chan string) {
	for {
		<-pingChan         // читаем ping
		pongChan <- "pong" // отправляем pong
	}
}

func main() {
	pingChan := make(chan string)
	pongChan := make(chan string)
	done := make(chan struct{})

	start := time.Now()
	const n = 10000000 // количество пар ping-pong

	go ping(pingChan, pongChan, done, n)
	go pong(pingChan, pongChan)

	<-done // ждем завершения

	elapsed := time.Since(start)

	fmt.Printf("Время: %v\n", elapsed)
	fmt.Printf("Сообщений (ping+pong): %d\n", n*2)
	fmt.Printf("Сообщений в секунду: %.2f\n", float64(n*2)/elapsed.Seconds())
}
