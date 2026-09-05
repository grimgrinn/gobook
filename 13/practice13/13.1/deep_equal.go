package main

import (
	"fmt"
	"math"
	"reflect"
)

const epsilon = 1e-9 // одна миллиардная

// Equal сравнивает два значения с учетом допуска для чисел с плавающей точкой
func Equal(x, y interface{}) bool {
	return equal(reflect.ValueOf(x), reflect.ValueOf(y))
}

func equal(x, y reflect.Value) bool {
	if x.Type() != y.Type() {
		return false
	}

	switch x.Kind() {
	case reflect.Invalid:
		return true

	case reflect.Bool:
		return x.Bool() == y.Bool()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return x.Int() == y.Int()

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return x.Uint() == y.Uint()

	case reflect.Float32, reflect.Float64:
		diff := math.Abs(x.Float() - y.Float())
		return diff < epsilon

	case reflect.String:
		return x.String() == y.String()

	case reflect.Ptr:
		if x.IsNil() && y.IsNil() {
			return true
		}
		if x.IsNil() || y.IsNil() {
			return false
		}
		return equal(x.Elem(), y.Elem())

	case reflect.Array:
		for i := 0; i < x.Len(); i++ {
			if !equal(x.Index(i), y.Index(i)) {
				return false
			}
		}
		return true

	case reflect.Slice:
		if x.IsNil() && y.IsNil() {
			return true
		}
		if x.IsNil() || y.IsNil() {
			return false
		}

		if x.Len() != y.Len() {
			return false
		}
		for i := 0; i < x.Len(); i++ {
			if !equal(x.Index(i), y.Index(i)) {
				return false
			}
		}
		return true

	case reflect.Struct:
		for i := 0; i < x.NumField(); i++ {
			if !equal(x.Field(i), y.Field(i)) {
				return false
			}
		}
		return true

	case reflect.Map:
		if x.IsNil() && y.IsNil() {
			return true
		}
		if x.IsNil() || y.IsNil() {
			return false
		}

		if len(x.MapKeys()) != len(y.MapKeys()) {
			return false
		}
		for _, key := range x.MapKeys() {
			xv := x.MapIndex(key)
			yv := y.MapIndex(key)
			if !yv.IsValid() || !equal(xv, yv) {
				return false
			}
		}
		return true

	case reflect.Interface:
		if x.IsNil() && y.IsNil() {
			return true
		}
		if x.IsNil() || y.IsNil() {
			return false
		}
		return equal(x.Elem(), y.Elem())

	default:
		return false
	}
}

type Person struct {
	Name  string
	Age   int
	Score float64
}

func main() {
	a := Person{
		Name:  "Alice",
		Age:   30,
		Score: 0.3000000001,
	}

	b := Person{
		Name:  "Alice",
		Age:   30,
		Score: 0.3000000002,
	}

	fmt.Println(Equal(a, b)) // true (разница < 1e-9)

	c := Person{
		Name:  "Alice",
		Age:   30,
		Score: 0.3000001,
	}

	fmt.Println(Equal(a, c)) // false (разница > 1e-9)
}
