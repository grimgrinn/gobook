package main

import (
	"fmt"
	"gobook/10/practice10/10.2/reader"
	_ "gobook/10/practice10/10.2/reader"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: archive-reader <file>")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	info, _ := f.Stat()
	ar, name, err := reader.Open(f, info.Size())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Формат архива: %s\n", name)

	for ar.Next() {
		file := ar.File()
		fmt.Printf("%s (%d байт)\n", file.Name(), file.Size())
	}
	if err := ar.Err(); err != nil {
		log.Fatal(err)
	}
}
