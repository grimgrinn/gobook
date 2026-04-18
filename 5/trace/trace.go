package main

import (
	"log"
	"time"
)

func bigSlowOperation() {
	defer trace("bigSLowOperation")() // не забываем о скобках
	// ... Длительная работа ,,,
	time.Sleep(10 * time.Second)
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("вход в %s", msg)
	return func() { log.Printf("выход из %s (%s)", msg, time.Since(start)) }
}

func main() {
	bigSlowOperation()
}
