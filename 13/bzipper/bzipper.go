// Bzipper читает входные данные, сжимвет их с помощью
// bzip2compress и записывает в стандартный вывод.
package main

import (
	"gobook/13/bzip"
	"io"
	"log"
	"os"
)

func main() {
	w := bzip.NewWriter(os.Stdout)
	if _, err := io.Copy(w, os.Stdin); err != nil {
		log.Fatalf("bzipper: %v\n", err)
	}

	if err := w.Close(); err != nil {
		log.Fatalf("bzipper: закрыт: %v\n", err)
	}
}
