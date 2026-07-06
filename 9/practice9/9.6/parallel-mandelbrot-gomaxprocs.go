package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	xmin, ymin, xmax, ymax = -2, -2, 2, 2
	width, height          = 1024, 1024
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

func renderRow(y int) []color.Color {
	row := make([]color.Color, width)
	for x := 0; x < width; x++ {
		xf := float64(x)/float64(width)*(xmax-xmin) + xmin
		yf := float64(y)/float64(height)*(ymax-ymin) + ymin
		z := complex(xf, yf)
		row[x] = mandelbrot(z)
	}
	return row
}

func main() {
	// измеряем для разных GOMAXPROCS
	for _, n := range []int{1, 2, 4, 8, 16, 32} {
		runtime.GOMAXPROCS(n)

		start := time.Now()

		img := image.NewRGBA(image.Rect(0, 0, width, height))
		var wg sync.WaitGroup
		wg.Add(height)

		for y := 0; y < height; y++ {
			go func(y int) {
				defer wg.Done()
				row := renderRow(y)
				for x, col := range row {
					img.Set(x, y, col)
				}
			}(y)
		}
		wg.Wait()

		elapsed := time.Since(start)

		fmt.Printf("GOMAXPROCS=%d: %v\n", n, elapsed)

		// Сохраняем картинку
		if n == 4 {
			f, _ := os.Create("mandelbrote.png")
			png.Encode(f, img)
			f.Close()
		}
	}

	fmt.Printf("Всего ядер CPU: %d\n", runtime.NumCPU())
}
