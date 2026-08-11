package main

import (
	"bytes"
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
		if v.IsNil() {
			buf.WriteString("null")
			return nil
		}
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
		fieldCount := 0
		for i := 0; i < v.NumField(); i++ {
			// Пропускаем нулевые поля
			if v.Field(i).IsZero() {
				continue
			}
			if fieldCount > 0 {
				buf.WriteByte(',')
			}
			fieldName := v.Type().Field(i).Name
			fmt.Fprintf(buf, "%q:", fieldName)
			if err := encodeJSON(buf, v.Field(i)); err != nil {
				return err
			}
			fieldCount++
		}
		buf.WriteByte('}')

	case reflect.Map:
		if v.IsNil() {
			buf.WriteString("null")
			return nil
		}
		buf.WriteByte('{')
		keys := v.MapKeys()
		keyCount := 0
		for _, key := range keys {
			// Пропускаем ключи с нулевым значением
			if key.IsZero() {
				continue
			}
			if keyCount > 0 {
				buf.WriteByte(',')
			}
			if err := encodeJSON(buf, key); err != nil {
				return err
			}
			buf.WriteByte(':')
			if err := encodeJSON(buf, v.MapIndex(key)); err != nil {
				return err
			}
			keyCount++
		}
		buf.WriteByte('}')

	default:
		return fmt.Errorf("неподдерживыемый тип для JSON: %s", v.Type())
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
	Name    string
	Age     int
	Email   string // пустое - пропускается
	Active  bool   // false - пропускается
	Tags    []string
	Ignored int // 0 - пропускается
}

func main() {
	p := Person{
		Name:   "Alice",
		Age:    30,
		Tags:   []string{"go", "programming"},
		Email:  "",    // пропускается
		Active: false, // пропускается
	}

	data, err := MarshalJSON(p)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(data))
}
