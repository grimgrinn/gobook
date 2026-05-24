// Clock1 является TCP-сервером, периодически выводящим время.
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	port := flag.String("port", "8000", "port number")
	flag.Parse()
	listener, err := net.Listen("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // Например, обрыв соединения
			continue
		}
		go handleConn(conn) // Обработка единственного подключения
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("\r15:04:05\n"))
		if err != nil {
			return // Например, отключение клиента
		}
		time.Sleep(1 * time.Second)
	}
}
