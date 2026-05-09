package main

import (
	"fmt"
	"gobook/7/eval"
	"io"
	"math"
	"net/http"
)

func parseAndCheck(s string) (eval.Expr, error) {
	if s == "" {
		return nil, fmt.Errorf("пустое выражение")
	}

	expr, err := eval.Parse(s)
	if err != nil {
		return nil, err
	}

	vars := make(map[eval.Var]bool)
	if err := expr.Check(vars); err != nil {
		return nil, err
	}
	for v := range vars {
		if v != "x" && v != "y" && v != "r" {
			return nil, fmt.Errorf("неизвестная переменна: %s", v)
		}
	}
	return expr, nil
}

func plot(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	expr, err := parseAndCheck(r.Form.Get("expr"))
	if err != nil {
		http.Error(w, "некорректное выражение: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	surface(w, func(x, y float64) float64 {
		r := math.Hypot(x, y) // Расстояние от (0,0)
		return expr.Eval(eval.Env{"x": x, "y": y, "r": r})
	})
}

const (
	width, height = 600, 320
	cells         = 100
	xyrange       = 30.0
	angle         = math.Pi / 6
)

func surface(w io.Writer, f func(x, y float64) float64) {

	var sin30, cos30 = math.Sin(angle), math.Cos(angle)
	xyscale := width / 2 / xyrange
	zscale := height * 0.4

	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j, xyrange, sin30, cos30, width, height, xyscale, zscale, f)
			bx, by := corner(i, j, xyrange, sin30, cos30, width, height, xyscale, zscale, f)
			cx, cy := corner(i, j+1, xyrange, sin30, cos30, width, height, xyscale, zscale, f)
			dx, dy := corner(i+1, j+1, xyrange, sin30, cos30, width, height, xyscale, zscale, f)

			fmt.Fprintf(w, "<polygon points='%g, %g %g, %g %g, %g %g, %g'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy)
		}
	}
	fmt.Fprintf(w, "</svg>")
}

func corner(i, j int, xyrange, sin30, cos30 float64, width, height int, xyscale, zscale float64, f func(x, y float64) float64) (sx, sy float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z := f(x, y)

	sx = float64(width)/2 + (x-y)*cos30*xyscale
	sy = float64(height)/2 + (x+y)*sin30*xyscale - z*zscale
	return
}

func main() {
	http.HandleFunc("/plot", plot)
	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
