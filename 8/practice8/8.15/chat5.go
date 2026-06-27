package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type client struct {
	name string
	ch   chan<- string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[string]client)

	for {
		select {
		case msg := <-messages:
			for _, cli := range clients {
				select {
				case cli.ch <- msg:
					// отправлено
				default:
					// клиент не готов - пропускаем
				}
			}
		case cli := <-entering:
			if len(clients) == 0 {
				cli.ch <- "В чате никого нет"
			} else {
				var names []string
				for name := range clients {
					names = append(names, name)
				}
				cli.ch <- "В чате: " + strings.Join(names, ", ")
			}
			clients[cli.name] = cli

		case cli := <-leaving:
			delete(clients, cli.name)
			close(cli.ch)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string, 10) // Исходящие сообщения клиентов
	go clientWriter(conn, ch)

	// Запрашиваем имя
	fmt.Fprintf(conn, "Введите ваше имя: ")
	input := bufio.NewScanner(conn)
	if !input.Scan() {
		conn.Close()
		return
	}
	who := input.Text()
	if who == "" {
		who = conn.RemoteAddr().String()
	}

	ch <- "Добро пожаловать, " + who
	messages <- who + " подключился"

	cl := client{name: who, ch: ch}
	entering <- cl

	timer := time.NewTimer(5 * time.Minute)

	go func() {
		<-timer.C
		fmt.Fprintf(conn, "Таймаут бездействия 5 минут\n")
		conn.Close()
	}()

	for input.Scan() {
		timer.Reset(5 * time.Minute)
		messages <- who + ": " + input.Text()
	}

	timer.Stop()
	leaving <- cl
	messages <- who + " отключился"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // Примечание: игнорируем ошибки сети
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
