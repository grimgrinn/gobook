package eval

import (
	"fmt"
	"math"
	"strings"
)

// Expr представляет арифметическое выражение
type Expr interface {
	// Eval возвращает значение данного Expr в среде env.
	Eval(env Env) float64
	// Check сообщает об ошибках в данном Expr и добавляет свои Vars.
	Check(vars map[Var]bool) error
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
		return b.x.Eval(env) - b.Eval(env)
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

// func TestEval(t *testing.T) {
// 	tests := []struct {
// 		expr string
// 		env  Env
// 		want string
// 	}{
// 		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
// 		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
// 		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
// 		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
// 		{"5 / 9 * (F - 32)", Env{"F": 32}, "0"},
// 		{"5 / 9 * (F - 32)", Env{"F": 212}, "100"},
// 	}
// 	var prevExpr string
// 	for _, test := range tests {
// 		// Выводит expr, только когда ого изменяется.
// 		if test.expr != prevExpr {
// 			fmt.Printf("\n%s\n", test.expr)
// 			prevExpr = test.expr
// 		}
// 		expr, err := Parse(test.expr)
// 		if err != nil {
// 			t.Error(err) // Ошибка анализа
// 			continue
// 		}
// 		got := fmt.Sprintf("%s.6g", expr.Eval(test.env))
// 		fmt.Printf("\t%v => %s\n", test.env, got)
// 		if got != test.want {
// 			t.Errorf("%s.Eval() в %v = %q, требуется %q\n", test.expr, test.env, got, test.want)
// 		}

// 	}
// }

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
	if !strings.ContainsRune("+-*", b.op) {
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

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}
