package sexpr_token

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"text/scanner"
)

type Token interface {
	String() string
}

type Symbol string

func (s Symbol) String() string { return string(s) }

type String string

func (s String) String() string { return string(s) }

type Int int64

func (i Int) String() string { return strconv.FormatInt(int64(i), 10) }

type StartList struct{}

func (s StartList) String() string { return "(" }

type EndList struct{}

func (e EndList) String() string { return ")" }

type TokenDecoder struct {
	r   *bufio.Reader
	lex *lexer
}

func NewTokenDecoder(r io.Reader) *TokenDecoder {
	return &TokenDecoder{r: bufio.NewReader(r)}
}

func (d *TokenDecoder) Token() (Token, error) {
	if d.lex == nil {
		d.lex = &lexer{scan: scanner.Scanner{}}
		d.lex.scan.Init(d.r)
		d.lex.next()
	}

	for d.lex.token == scanner.EOF {
		return nil, io.EOF
	}

	switch d.lex.token {
	case scanner.EOF:
		return nil, io.EOF

	case scanner.Ident:
		s := d.lex.text()
		d.lex.next()
		return Symbol(s), nil

	case scanner.String:
		s := d.lex.text()
		s = s[1 : len(s)-1]
		d.lex.next()
		return String(s), nil

	case scanner.Int:
		s := d.lex.text()
		d.lex.next()
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		return Int(n), nil

	case '(':
		d.lex.next()
		return StartList{}, nil

	case ')':
		d.lex.next()
		return EndList{}, nil

	default:
		return nil, fmt.Errorf("неожиданный токен: %q", d.lex.text())
	}
}

func (d *TokenDecoder) Decode(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("Decode требует указатель на значение")
	}

	tok, err := d.Token()
	if err != nil {
		return err
	}
	if _, ok := tok.(StartList); !ok {
		return fmt.Errorf("ожидается начало списка, получено %T", tok)
	}

	return decodeTokenList(d, val.Elem())
}

func decodeTokenList(d *TokenDecoder, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Struct:
		for {
			tok, err := d.Token()
			if err != nil {
				return err
			}
			if _, ok := tok.(EndList); ok {
				return nil
			}

			if _, ok := tok.(StartList); !ok {
				return fmt.Errorf("ожидается начало списка, получено %T", tok)
			}

			nameTok, err := d.Token()
			if err != nil {
				return err
			}

			name, ok := nameTok.(Symbol)
			if !ok {
				return fmt.Errorf("ожидается имя поля (Symbol), получено %T", nameTok)
			}

			field := v.FieldByName(string(name))
			if !field.IsValid() {
				return fmt.Errorf("поле %s не найдено", name)
			}

			if err := decodeTokenValue(d, field); err != nil {
				return err
			}

			endTok, err := d.Token()
			if err != nil {
				return err
			}
			if _, ok := endTok.(EndList); !ok {
				return fmt.Errorf("ожидается конец списка, получено %T", endTok)
			}
		}
	case reflect.Slice:
		for {
			tok, err := d.Token()
			if err != nil {
				return err
			}
			if _, ok := tok.(EndList); ok {
				return nil
			}
			elem := reflect.New(v.Type().Elem()).Elem()
			if err := decodeTokenValue(d, elem); err != nil {
				return err
			}
			v.Set(reflect.Append(v, elem))
		}

	default:
		return fmt.Errorf("неподдерживаемый тип для списка: %s", v.Kind())
	}
}

func decodeTokenValue(d *TokenDecoder, v reflect.Value) error {
	tok, err := d.Token()
	if err != nil {
		return err
	}

	switch t := tok.(type) {
	case StartList:
		return decodeTokenList(d, v)

	case Symbol:
		s := string(t)
		if s == "t" {
			v.SetBool(true)
			return nil
		}
		if s == "nil" {
			v.SetZero()
			return nil
		}
		if v.Kind() == reflect.String {
			v.SetString(s)
			return nil
		}
		if v.Kind() == reflect.Int || v.Kind() == reflect.Int64 {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return err
			}
			v.SetInt(n)
			return nil
		}

		return fmt.Errorf("невозможно преобразовать Symbol %q в %s", s, v.Kind())

	case String:
		v.SetString(string(t))
		return nil

	case Int:
		if v.Kind() == reflect.Int || v.Kind() == reflect.Int64 {
			v.SetInt(int64(t))
			return nil
		}
		return fmt.Errorf("невозможно преобразовать Int %d в %s", t, v.Kind())

	default:
		return fmt.Errorf("неожиданный токен при декодировании значения: %T", tok)
	}
}

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }
