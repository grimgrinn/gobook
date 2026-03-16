package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка чтения: %v\n", err)
		os.Exit(1)
	}

	shaSize := flag.Int("sha", 256, "размер SHA (256, 384, 512)")
	flag.Parse()

	switch *shaSize {
	case 256:
		fmt.Printf("%x", sha256.Sum256(data))
	case 384:
		fmt.Printf("%x", sha512.Sum384(data))
	case 512:
		fmt.Printf("%x", sha512.Sum512(data))
	default:
		fmt.Fprintf(os.Stderr, "неподдерживаемый размер: %d\n", *shaSize)
		os.Exit(1)
	}
}
