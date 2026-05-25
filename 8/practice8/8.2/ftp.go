package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type client struct {
	conn net.Conn
	cwd  string // текущая рабочая директория
}

func (c *client) writeln(s string) {
	fmt.Fprintln(c.conn, s)
}

func (c *client) handle() {
	defer c.conn.Close()

	for {
		netData, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			log.Println("client disconnected")
			return
		}

		cmd := strings.TrimSpace(strings.ToLower(netData))
		args := strings.Fields(cmd)
		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "ls":
			entries, err := os.ReadDir(c.cwd)
			if err != nil {
				c.writeln("error: " + err.Error())
				continue
			}
			for _, entry := range entries {
				c.writeln(entry.Name())
			}
			c.writeln("OK") // маркер конца списка

		case "cd":
			if len(args) < 2 {
				c.writeln("usage: cd <dir>")
				continue
			}
			newPath := filepath.Join(c.cwd, args[1])
			info, err := os.Stat(newPath)
			if err != nil || !info.IsDir() {
				c.writeln("error: invalid directory")
				continue
			}
			c.cwd = newPath
			c.writeln("OK")

		case "get":
			if len(args) < 2 {
				c.writeln("usage: get <file>")
				continue
			}
			filePath := filepath.Join(c.cwd, args[1])
			file, err := os.Open(filePath)
			if err != nil {
				c.writeln("error: " + err.Error())
				continue
			}
			defer file.Close()
			_, err = io.Copy(c.conn, file)
			if err != nil {
				log.Println("error sending file:", err)
			}
			c.writeln("OK")

		case "close", "quit", "exit":
			c.writeln("bye")
			return
		default:
			c.writeln("unknow command: " + args[0])
		}
	}
}

func main() {
	port := flag.String("port", "8021", "ftp port")
	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("FTP server listening on port %s", *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go func(c net.Conn) {
			// Получаем текущую директорию сервера (можно задать отдельную папку)
			startDir, _ := os.Getwd()
			client := &client{
				conn: c,
				cwd:  startDir,
			}
			client.writeln("Welcome to Go FTP")
			client.handle()
		}(conn)
	}
}
