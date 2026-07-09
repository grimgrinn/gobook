package main

import (
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

var format = flag.String("fmt", "jpeg", "выходной формат: jpeg, png, gif")

func main() {
	flag.Parse()

	if err := convert(os.Stdin, os.Stdout, *format); err != nil {
		fmt.Fprintf(os.Stderr, "convert: %v\n", err)
		os.Exit(1)
	}
}

func convert(in io.Reader, out io.Writer, format string) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return fmt.Errorf("декодирование: %v", err)
	}
	fmt.Fprintf(os.Stderr, "Входной формат: %s\n", kind)

	switch format {
	case "jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	case "png":
		return png.Encode(out, img)
	case "gif":
		return gif.Encode(out, img, &gif.Options{NumColors: 256})
	default:
		return fmt.Errorf("неподдерживаемый формат: %s", format)
	}
}
