package main

import "fmt"

type Point struct {
	x, y int
}

type Circle struct {
	Point
	Radius int
}

type Wheel struct {
	Circle
	Spokes int
}

func main() {
	//w = Wheel{Circle{Point{8,8}, 5}, 20}
	w := Wheel{
		Circle: Circle{
			Point:  Point{x: 8, y: 8},
			Radius: 5,
		},
		Spokes: 20, // завершающая запятая необдима
	}
	fmt.Printf("%#v\n", w)

	w.x = 42

	fmt.Printf("%#v\n", w)
}
