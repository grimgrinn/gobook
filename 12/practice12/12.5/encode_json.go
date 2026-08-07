package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

func encodeJSON(buf *bytes.Buffer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("null")

	case reflect.Bool:
		if v.Bool() {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%g", v.Float())

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		if v.IsNil() {
			buf.WriteString("null")
			return nil
		}
		return encodeJSON(buf, v.Elem())

	case reflect.Array, reflect.Slice:
		buf.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := encodeJSON(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(']')

	case reflect.Struct:
		buf.WriteByte('{')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(',')
			}
			fieldName := v.Type().Field(i).Name
			fmt.Fprintf(buf, "%q:", fieldName)
			if err := encodeJSON(buf, v.Field(i)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	case reflect.Map:
		buf.WriteByte('{')
		keys := v.MapKeys()
		for i, key := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := encodeJSON(buf, key); err != nil {
				return err
			}
			buf.WriteByte(':')
			if err := encodeJSON(buf, v.MapIndex(key)); err != nil {
				return err
			}
		}
		buf.WriteByte('}')

	default:
		return fmt.Errorf("неподерживаемый тип для JSON: %s", v.Type())
	}
	return nil
}

func MarshalJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encodeJSON(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Person struct {
	Name string
	Age  int
	Tags []string
}

func main() {
	p := Person{
		Name: "Alice",
		Age:  30,
		Tags: []string{"go", "programming"},
	}

	data, err := MarshalJSON(p)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}
	fmt.Println("Our JSON:", string(data))

	var decoded Person
	if err := json.Unmarshal(data, &decoded); err != nil {
		fmt.Println("Unmarshal error:", err)
		return
	}
	fmt.Printf("Decoded: %+v\n", decoded)
}
