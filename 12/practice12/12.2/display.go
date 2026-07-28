package main

import (
	"fmt"
	"reflect"
	"strings"
)

func formatKey(key reflect.Value) string {
	switch key.Kind() {
	case reflect.Struct:
		return formatStruct(key)
	case reflect.Array:
		return formatArray(key)
	case reflect.Slice:
		return formatSlice(key)
	case reflect.Map:
		return formatMap(key)
	default:
		return fmt.Sprintf("%v", key.Interface())
	}
}

func formatStruct(v reflect.Value) string {
	var parts []string
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		parts = append(parts, fmt.Sprintf("%s:%v", field.Name, v.Field(i).Interface()))
	}
	return "{" + strings.Join(parts, " ") + "}"
}

func formatArray(v reflect.Value) string {
	var parts []string
	for i := 0; i < v.Len(); i++ {
		parts = append(parts, fmt.Sprintf("%v", v.Index(i).Interface()))
	}
	return "[" + strings.Join(parts, " ") + "]"
}

func formatSlice(v reflect.Value) string {
	var parts []string
	for i := 0; i < v.Len(); i++ {
		parts = append(parts, fmt.Sprintf("%v", v.Index(i).Interface()))
	}
	return "[" + strings.Join(parts, " ") + "]"
}

func formatMap(v reflect.Value) string {
	var parts []string
	for _, key := range v.MapKeys() {
		parts = append(parts, fmt.Sprintf("%s:%v", formatKey(key), v.MapIndex(key).Interface()))
	}
	return "{" + strings.Join(parts, " ") + "}"
}

const maxDepth = 100

func Display(name string, x interface{}) {
	fmt.Printf("Display %s (%T):\n", name, x)
	display(name, reflect.ValueOf(x), 0)
}

func display(path string, v reflect.Value, depth int) {
	if depth > maxDepth {
		fmt.Printf("%s = ... (глубина превышена)\n", path)
		return
	}

	switch v.Kind() {
	case reflect.Invalid:
		fmt.Printf("%s = invalid\n", path)
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i), depth+1)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fieldPath := fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name)
			display(fieldPath, v.Field(i), depth+1)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatKey(key)), v.MapIndex(key), depth+1)
		}
	case reflect.Ptr:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("%s", path), v.Elem(), depth+1)
		}
	case reflect.Interface:
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("%s", path), v.Elem(), depth+1)
		}
	default:
		fmt.Printf("%s = %v\n", path, v)
	}

}

type Node struct {
	Value int
	Next  *Node
}

func main() {
	a := &Node{Value: 1}
	b := &Node{Value: 2}
	a.Next = b
	b.Next = a // цикл!

	Display("cycle", a)
}
