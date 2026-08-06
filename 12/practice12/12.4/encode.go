package main

import (
	"bytes"
	"fmt"
	"reflect"
)

func encodeWithIndent(buf *bytes.Buffer, v reflect.Value, indent string) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")

	case reflect.Bool:
		if v.Bool() {
			buf.WriteString("t")
		} else {
			buf.WriteString("nil")
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprintf(buf, "%d", v.Uint())

	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%g", v.Float())

	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		fmt.Fprintf(buf, "#C(%g %g)", real(c), imag(c))

	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())

	case reflect.Ptr:
		if v.IsNil() {
			buf.WriteString("nil")
			return nil
		}
		return encodeWithIndent(buf, v.Elem(), indent)

	case reflect.Array, reflect.Slice:
		buf.WriteString("(\n")
		nextIndext := indent + " "
		for i := 0; i < v.Len(); i++ {
			buf.WriteString(nextIndext)
			if err := encodeWithIndent(buf, v.Index(i), nextIndext); err != nil {
				return err
			}
			if i < v.Len()-1 {
				buf.WriteByte('\n')
			}
		}
		buf.WriteString("\n" + indent + ")")

	case reflect.Struct:
		buf.WriteString("(\n")
		nextIndent := indent + " "
		for i := 0; i < v.NumField(); i++ {
			fieldName := v.Type().Field(i).Name
			buf.WriteString(nextIndent)
			fmt.Fprintf(buf, "(%s ", fieldName)
			if err := encodeWithIndent(buf, v.Field(i), nextIndent); err != nil {
				return err
			}
			buf.WriteByte(')')
			if i < v.NumField()-1 {
				buf.WriteByte('\n')
			}
		}
		buf.WriteString("\n" + indent + ")")

	case reflect.Map:
		buf.WriteString("(\n")
		nextIndent := indent + " "
		keys := v.MapKeys()
		for i, key := range keys {
			buf.WriteString(nextIndent)
			buf.WriteByte('(')
			if err := encodeWithIndent(buf, key, nextIndent); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encodeWithIndent(buf, v.MapIndex(key), nextIndent); err != nil {
				return err
			}
			buf.WriteByte(')')
			if i < len(keys)-1 {
				buf.WriteByte('\n')
			}
		}
		buf.WriteString("\n" + indent + ")")

	case reflect.Interface:
		if v.IsNil() {
			buf.WriteString("nil")
			return nil
		}
		inner := v.Elem()
		typeName := inner.Type().String()
		fmt.Fprintf(buf, "(%s ", typeName)
		if err := encodeWithIndent(buf, inner, indent+" "); err != nil {
			return err
		}
		buf.WriteByte(')')

	default:
		return fmt.Errorf("неподдерживаемый тип: %s", v.Type())
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encodeWithIndent(&buf, reflect.ValueOf(v), ""); err != nil {
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

	data, err := Marshal(p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(data))
}
