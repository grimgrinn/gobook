package sexpr

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
)

var typeRegistry = make(map[string]reflect.Type)

func RegisterType(name string, typ reflect.Type) {
	typeRegistry[name] = typ
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
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
		return encode(buf, v.Elem())

	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			if err := encode(buf, v.Index(i)); err != nil {
				return err
			}
		}
		buf.WriteByte(')')

	case reflect.Struct:
		buf.WriteByte('(')
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteByte(' ')
			}
			name := v.Type().Field(i).Name
			fmt.Fprintf(buf, "(%s ", name)
			if err := encode(buf, v.Field(i)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case reflect.Map:
		buf.WriteByte('(')
		keys := v.MapKeys()
		for i, key := range keys {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteByte('(')
			if err := encode(buf, key); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key)); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')

	case reflect.Interface:
		if v.IsNil() {
			buf.WriteString("nil")
			return nil
		}
		inner := v.Elem()
		typeName := inner.Type().String()
		fmt.Fprintf(buf, "(%s ", typeName)
		if err := encode(buf, inner); err != nil {
			return err
		}
		buf.WriteByte(')')

	default:
		return fmt.Errorf("неподдерживыемый тип: %s", v.Type())
	}
	return nil
}

func Unmarshal(data []byte, v interface{}) error {
	dec := NewDecoder(bytes.NewReader(data))
	return dec.Decoder(v)
}

type Decoder struct {
	r   *bufio.Reader
	lex *lexer
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

func (d *Decoder) Decoder(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("Decode требует указатель на значеие")
	}

	d.lex = &lexer{scan: scanner.Scanner{}}
	d.lex.scan.Init(d.r)
	d.lex.next()

	if d.lex.token == scanner.EOF {
		return io.EOF
	}

	expr, err := parse(d.lex)
	if err != nil {
		return err
	}

	return decodeValue(expr, val.Elem())
}

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }
func (lex *lexer) consume(want rune) {
	if lex.token != want {
		panic(fmt.Sprintf("получечегл %q, требуются %q", lex.text(), want))
	}
	lex.next()
}

func parse(lex *lexer) (interface{}, error) {
	switch lex.token {
	case scanner.EOF:
		return nil, io.EOF
	case scanner.Ident:
		return parseIdent(lex)
	case scanner.String:
		return parseString(lex)
	case scanner.Int, scanner.Float:
		return parseNumber(lex)
	case '(':
		return parseList(lex)
	default:
		return nil, fmt.Errorf("неожиданный токен: %q", lex.text())
	}
}

func parseIdent(lex *lexer) (interface{}, error) {
	s := lex.text()
	lex.next()
	return s, nil
}

func parseString(lex *lexer) (interface{}, error) {
	s := lex.text()
	s = s[1 : len(s)-1]
	lex.next()
	return s, nil
}

func parseNumber(lex *lexer) (interface{}, error) {
	s := lex.text()
	lex.next()
	if strings.Contains(s, ".") || strings.Contains(s, "e") {
		return strconv.ParseFloat(s, 64)
	}
	return strconv.ParseInt(s, 10, 64)
}

func parseList(lex *lexer) ([]interface{}, error) {
	lex.consume('(')
	var list []interface{}
	for {
		switch lex.token {
		case ')':
			lex.next()
			return list, nil
		case scanner.EOF:
			return nil, fmt.Errorf("незакрытый список")

		default:
			item, err := parse(lex)
			if err != nil {
				return nil, err
			}
			list = append(list, item)
		}
	}
}

func decodeValue(expr interface{}, v reflect.Value) error {
	switch val := expr.(type) {
	case nil:
		v.SetZero()
		return nil

	case string:
		if val == "t" {
			v.SetBool(true)
			return nil
		}
		if val == "nil" {
			v.SetZero()
			return nil
		}

		if v.Kind() == reflect.Int || v.Kind() == reflect.Int64 {
			n, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				v.SetInt(n)
				return nil
			}
		}

		if v.Kind() == reflect.Interface {
			if _, ok := typeRegistry[val]; ok {

				return fmt.Errorf("интерейс: ожидается список для типа %s", val)

			}
		}

		v.SetString(val)
		return nil

	case int64:
		if v.Kind() == reflect.Int || v.Kind() == reflect.Int64 {
			v.SetInt(val)
			return nil
		}
		if v.Kind() == reflect.Float64 || v.Kind() == reflect.Float32 {
			v.SetFloat(float64(val))
			return nil
		}
		return fmt.Errorf("невозможно преобразовать %d в %s", val, v.Kind())

	case float64:
		if v.Kind() == reflect.Float64 || v.Kind() == reflect.Float32 {
			v.SetFloat(val)
			return nil
		}
		if v.Kind() == reflect.Int || v.Kind() == reflect.Int64 {
			v.SetInt(int64(val))
			return nil
		}
		return fmt.Errorf("невозможно преобразовать %f в %s", val, v.Kind())

	case []interface{}:
		return decodeList(val, v)

	default:
		return fmt.Errorf("неизвестный тип: %T", expr)
	}
}

func decodeList(list []interface{}, v reflect.Value) error {
	if v.Kind() == reflect.Interface {
		if len(list) != 2 {
			return fmt.Errorf("интерфейс должен быть закодирован как (тип значение)")
		}
		typeName, ok := list[0].(string)
		if !ok {
			return fmt.Errorf("имя типа должно быть строкой")
		}
		typ, ok := typeRegistry[typeName]
		if !ok {
			return decodeStruct(list, v)
		}
		newVal := reflect.New(typ).Elem()
		if err := decodeValue(list[1], newVal); err != nil {
			return err
		}
		v.Set(newVal)
		return nil
	}

	switch v.Kind() {
	case reflect.Slice:
		for _, elem := range list {
			newVal := reflect.New(v.Type().Elem()).Elem()
			if err := decodeValue(elem, newVal); err != nil {
				return err
			}
			v.Set(reflect.Append(v, newVal))
		}
		return nil

	case reflect.Struct:
		return decodeStruct(list, v)

	default:
		return fmt.Errorf("неподдерживаемый тип для списка: %s", v.Kind())
	}
}

func decodeStruct(list []interface{}, v reflect.Value) error {
	for _, elem := range list {
		pair, ok := elem.([]interface{})
		if !ok || len(pair) != 2 {
			return fmt.Errorf("ожидается пара (имя значение)")

		}
		name, ok := pair[0].(string)
		if !ok {
			return fmt.Errorf("имя поля должно быть строкой")

		}
		field := v.FieldByName(name)
		if !field.IsValid() {
			return fmt.Errorf("поле %s не найдено", name)

		}
		if err := decodeValue(pair[1], field); err != nil {
			return err
		}
	}

	return nil
}
