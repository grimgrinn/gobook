package main

import (
	"fmt"
	"gobook/12/params"
)

type SearchData struct {
	Labels     []string `http:"labels"`
	MaxResults int      `http:"max"`
	Exact      bool     `http:"x"`
}

func main() {
	data := SearchData{
		Labels:     []string{"go", "programming"},
		MaxResults: 5,
		Exact:      true,
	}

	query, err := params.Pack(&data)
	if err != nil {
		fmt.Println("Errors:", err)
		return
	}

	fmt.Println("URL параметры:", query)
}
