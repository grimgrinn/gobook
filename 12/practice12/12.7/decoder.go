package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Decoder - потоковый декодер S-выражений
type Decoder struct {
	r *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

// Decode читает S-выражения и заполняет v
func (d *Decoder) Decode(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("Decode requres a non-nil pointer")
	}

	expr, err := parse(d.r)
	if err != nil {
		return err
	}

	return decodeValue(expr, val.Elem())
}

type token struct {
	typ tokenType
	val string
}

type tokenType int

const (
	tokenLParen tokenType = iota
	tokenParen
	tokenSymbol
	tokenString
	tokenNumber
	tokenEOF
	tokenError
)

func parse(r *bufio.Reader) (interface{}, error) {
	skipWhitespace(r)
	ch, err := peekChar(r)
	if err != nil {
		return nil, err
	}

	switch ch {
	case '(':
		return parseList(r)
	case '"':
		return parseString(r)
	default:
		if isDigit(ch) || ch == '-' {
			return parseNumber(r)
		}
		return parseSymbol(r)
	}
}

func parseList(r *bufio.Reader) ([]interface{}, error) {
	r.ReadByte()

	var list []interface{}
	for {
		skipWhitespace(r)
		ch, err := peekChar(r)
		if err != nil {
			return nil, err
		}
		if ch == ')' {
			r.ReadByte()
			return list, nil
		}
		elem, err := parse(r)
		if err != nil {
			return nil, err
		}
		list = append(list, elem)
	}
}

func parseString(r *bufio.Reader) (string, error) {
	r.ReadByte()
	var buf strings.Builder

	for {
		ch, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if ch == '"' {
			break
		}
		buf.WriteByte(ch)
	}
	return buf.String(), nil
}

func parseNumber(r *bufio.Reader) (interface{}, error) {
	var buf strings.Builder

	for {
		ch, err := peekChar(r)
		if err != nil {
			break
		}

		if !isDigit(ch) && ch != '.' && ch != '-' {
			break
		}
		r.ReadByte()
		buf.WriteByte(ch)
	}
	s := buf.String()
	if strings.Contains(s, ",") {
		return strconv.ParseFloat(s, 64)
	}
	return strconv.ParseInt(s, 10, 64)
}

func parseSymbol(r *bufio.Reader) (string, error) {
	var buf strings.Builder
	for {
		ch, err := peekChar(r)
		if err != nil {
			break
		}
		if isWhitespace(ch) || ch == '(' || ch == ')' || ch == '"' {
			break
		}
		r.ReadByte()
		buf.WriteByte(ch)
	}
	return buf.String(), nil
}

func peekChar(r *bufio.Reader) (byte, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	r.UnreadByte()
	return b, nil
}

func skipWhitespace(r *bufio.Reader) {
	for {
		ch, err := peekChar(r)
		if err != nil || !isWhitespace(ch) {
			break
		}
		r.ReadByte()
	}
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch == 'g'
}

func decodeValue(expr interface{}, v reflect.Value) error {
	switch val := expr.(type) {
	case nil:
		return nil
	case string:
		if v.Kind() == reflect.Int || v.Kind() == reflect.Int64 {
			num, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				v.SetInt(num)
				return nil
			}
		}
		if v.Kind() == reflect.Float64 || v.Kind() == reflect.Float32 {
			num, err := strconv.ParseFloat(val, 64)
			if err == nil {
				v.SetFloat(num)
				return nil
			}
		}
		if val == "t" {
			v.SetBool(true)
			return nil
		}
		if val == "nil" {
			v.SetZero()
			return nil
		}
		v.SetString(val)
		return nil

	case int64:
		v.SetInt(val)
		return nil
	case float64:
		v.SetFloat(val)
		return nil
	case []interface{}:
		return decodeList(val, v)
	default:
		return fmt.Errorf("неизвестный тип: %T", expr)
	}
}

func decodeList(list []interface{}, v reflect.Value) error {
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
		for _, elem := range list {
			pair, ok := elem.([]interface{})
			if !ok || len(pair) != 2 {
				return fmt.Errorf("ожидается пара (имя значение)")
			}
			name, ok := pair[0].(string)
			if !ok {
				return fmt.Errorf("Имя поля должно быть строкой")
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

	default:
		return fmt.Errorf("неподдерживаемый тип для списка: %s", v.Kind())
	}
}

type Person struct {
	Name string
	Age  int
	Tags []string
}

func main() {
	s := `( (Name "Alice") (Age 30) (Tags ("go" "programming")) )`
	r := bytes.NewReader([]byte(s))

	var p Person
	dec := NewDecoder(r)

	if err := dec.Decode(&p); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Decoded: %+v\n", p)
}
