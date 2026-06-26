package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
				cli.ch <- msg
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
	ch := make(chan string) // Исходящие сообщения клиентов
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "Вы " + who
	messages <- who + " подключился"

	cl := client{name: who, ch: ch}
	entering <- cl

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}

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
