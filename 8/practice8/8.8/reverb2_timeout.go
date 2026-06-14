package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for {
		// Создаем канал, который сработает через 10 секунд
		timer := time.After(10 * time.Second)

		// Канал для сигнала о завершении сканирования
		scanDone := make(chan bool)

		// Запускаем сканирование в отдельной горутине
		go func() {
			scanDone <- input.Scan()
		}()

		select {
		case ok := <-scanDone:
			if !ok {
				// Ошибка или EOF
				c.Close()
				return
			}
			// Получили строку, обрабатываем
			go echo(c, input.Text(), 1*time.Second)
			// Продолжаем цикл (таймер сбросится на следующей итерации)
		case <-timer:
			// Тайм-аут: клиент ничего не отправлял 10 секунд
			fmt.Fprintf(c, "тайм-аут 10 секунд, закрываю соединение\n")
			c.Close()
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
