package main

import (
	"bytes"
	"fmt"
	"gobook/12/practice12/12.9/sexpr_token"
)

type Person struct {
	Name string
	Age  int
	Tags []string
}

func main() {
	s := `((Name "Alice")) (Age 30) (Tags ("go" "programming")))`
	r := bytes.NewReader([]byte(s))

	var p Person
	dec := sexpr_token.NewTokenDecoder(r)
	if err := dec.Decode(&p); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Decode: %+v\n", p)
}
