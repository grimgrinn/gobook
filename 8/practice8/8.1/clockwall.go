package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type clock struct {
	city string
	time string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: clockwall City=host:port [City2=hos2:port2 ...]")
	}

	clocks := make([]clock, 0, len(os.Args)-1)
	ch := make(chan clock)

	for _, arg := range os.Args[1:] {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			log.Fatalf("invalid format: %s (expected City=host:port)", arg)
		}
		city := parts[0]
		addr := parts[1]

		go func(city, addr string) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			for {
				var buf [64]byte
				n, err := conn.Read(buf[:])
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Print(err)
					break
				}
				ch <- clock{city, string(buf[:n])}
			}
		}(city, addr)

		clocks = append(clocks, clock{city, ""})
	}

	// Основной цикл обновления вывода
	for {
		c := <-ch
		// Обновляем время в слайсе clocks
		for i := range clocks {
			if clocks[i].city == c.city {
				clocks[i].time = strings.TrimSpace(c.time)
			}
		}
		// Очищаем экран и выводим таблицу
		fmt.Printf("\r\033[K") // очистка экрана (ANSI)
		for _, cl := range clocks {
			fmt.Printf("%s: %s\t", cl.city, cl.time)
		}
	}
}
