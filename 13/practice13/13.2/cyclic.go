package main

import (
	"fmt"
	"reflect"
)

// IsCyclic проверяет, содержит ли структура данных циклические ссылки.
func IsCyclic(x interface{}) bool {
	seen := make(map[uintptr]bool)
	return isCyclic(reflect.ValueOf(x), seen)
}

func isCyclic(v reflect.Value, seen map[uintptr]bool) bool {
	if !v.IsValid() {
		return false
	}

	if v.CanAddr() {
		switch v.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Func, reflect.Slice, reflect.UnsafePointer:
			ptr := v.Pointer()
			if seen[ptr] {
				return true
			}
			seen[ptr] = true
			defer func() { delete(seen, ptr) }()
		}
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return false
		}
		return isCyclic(v.Elem(), seen)

	case reflect.Interface:
		if v.IsNil() {
			return false
		}
		return isCyclic(v.Elem(), seen)

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if isCyclic(v.Field(i), seen) {
				return true
			}
		}
		return false

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if isCyclic(v.Index(i), seen) {
				return true
			}
		}
		return false

	case reflect.Map:
		for _, key := range v.MapKeys() {
			if isCyclic(key, seen) {
				return true
			}
			if isCyclic(v.MapIndex(key), seen) {
				return true
			}
		}
		return false

	default:
		return false
	}
}

type Node struct {
	Value int
	Next  *Node
}

type CyclicStruct struct {
	Self *CyclicStruct
}

func main() {
	//  обычный связный список (без цикла)
	a := &Node{Value: 1}
	b := &Node{Value: 2}
	c := &Node{Value: 3}
	a.Next = b
	b.Next = c

	fmt.Println("Cписок без цикла:", IsCyclic(a)) // false

	a.Next = b
	b.Next = c
	c.Next = a // зациклили

	fmt.Println("Циклический список:", IsCyclic(a)) // true

	// структура, указывающая на себя
	s := &CyclicStruct{}
	s.Self = s

	fmt.Println("Самоуказывающая структура:", IsCyclic(s)) // true

	// обычный cлайс (без цикла)
	slice := []int{1, 2, 3}
	fmt.Println("Обычный слайс:", IsCyclic(slice)) // false

	// слайс, содержащий сам себя (цикл без указателья)
	ptrSlice := []interface{}{}
	ptrSlice = append(ptrSlice, &ptrSlice)
	fmt.Println("Слайс с собой:", IsCyclic(ptrSlice)) // true
}
