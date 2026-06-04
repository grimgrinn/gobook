package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"
)

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 15

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return color.Gray{255 - contrast*n}
		}
	}
	return color.Black
}

func renderRow(y int, width int, xmin, xmax, ymin, ymax float64) []color.Color {
	row := make([]color.Color, width)
	for x := 0; x < width; x++ {
		xf := float64(x)/float64(width)*(xmax-xmin) + xmin
		yf := float64(y)/float64(height)*(ymax-ymin) + ymin
		z := complex(xf, yf)
		row[x] = mandelbrot(z)
	}
	return row
}

const (
	xmin, ymin, xmax, ymax = -2, -2, 2, 2
	width, height          = 1024, 1024
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var wg sync.WaitGroup

	for py := 0; py < height; py++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			row := renderRow(y, width, xmin, xmax, ymin, ymax)
			for x, col := range row {
				img.Set(x, y, col)
			}
		}(py)
	}
	wg.Wait()
	png.Encode(os.Stdout, img)
}
