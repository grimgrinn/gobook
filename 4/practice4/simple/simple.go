package main

import "fmt"

type Currency int

const (
	USD Currency = iota
	EUR
	GBP
	RUR
)

func main() {
	symbol := [...]string{USD: "$", EUR: "€", GBP: "£", RUR: "₽"}
	fmt.Println(RUR, symbol[RUR]) // "3 ₽"
}
