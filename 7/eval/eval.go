package eval

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Expr представляет арифметическое выражение
type Expr interface {
	// Eval возвращает значение данного Expr в среде env.
	Eval(env Env) float64
	// Check сообщает об ошибках в данном Expr и добавляет свои Vars.
	Check(vars map[Var]bool) error

	String() string
}

// Var определяет переменную, например x.
type Var string

// literal представляет собой числовуб константу, например 3.141.
type literal float64

// unary представляет выражение с унарным оператором, например -x.
type unary struct {
	op rune // '+' или '-'
	x  Expr
}

// binary представляет выражение с бинарным оператором, например x+y.
type binary struct {
	op   rune // '+', '-', '*' или '/'
	x, y Expr
}

// call представляет выражаение вызова функции, например sin(x).
type call struct {
	fn   string // одна из "pow", "sin", "sqrt"
	args []Expr
}

// min представляет выражение вычисления минимального значения.
type min struct {
	args []Expr
}

type Env map[Var]float64

func (v Var) Eval(env Env) float64 {
	return env[v]
}

func (l literal) Eval(_ Env) float64 { // типо просто возвращает l O_o
	return float64(l)
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("неподдерживаемый унарный оператор: %q", u.op))
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("неподдерживаемый бинарный оператор: %q", b.op))
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("неподдерживаемый вызов функции: %s", c.fn))
}

func (m min) Eval(env Env) float64 {
	if len(m.args) == 0 {
		return 0
	}
	minVal := m.args[0].Eval(env)
	for _, arg := range m.args[1:] {
		v := arg.Eval(env)
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

func (literal) Check(vars map[Var]bool) error {
	return nil
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("некорректный унарный оператор %q", u.op)
	}
	return u.x.Check(vars)
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf("некорректный бинарный оператор %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("неизвестная функция %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("вызов %s имеет %d вместо %d аргументов", c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

func (m min) Check(vars map[Var]bool) error {
	if len(m.args) == 0 {
		return fmt.Errorf("min: отсутствуют аргументы")
	}
	for _, arg := range m.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

func (v Var) String() string {
	return string(v)
}

func (l literal) String() string {
	return strconv.FormatFloat(float64(l), 'g', -1, 64)
}

func (u unary) String() string {
	return string(u.op) + u.x.String()
}

func (b binary) String() string {
	return b.x.String() + string(b.op) + b.y.String()
}

func (c call) String() string {
	args := make([]string, len(c.args))
	for i, arg := range c.args {
		args[i] = arg.String()
	}
	return c.fn + "(" + strings.Join(args, ", ") + ")"
}

func (m min) String() string {
	args := make([]string, len(m.args))
	for i, arg := range m.args {
		args[i] = arg.String()
	}
	return "min(" + strings.Join(args, ", ") + ")"
}
