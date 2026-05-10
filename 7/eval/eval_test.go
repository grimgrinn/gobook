package eval

import (
	"fmt"
	"math"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
		{"5 / 9 * (F - 32)", Env{"F": 32}, "0"},
		{"5 / 9 * (F - 32)", Env{"F": 212}, "100"},
	}
	var prevExpr string
	for _, test := range tests {
		// Выводит expr, только когда ого изменяется.
		if test.expr != prevExpr {
			fmt.Printf("\n%s\n", test.expr)
			prevExpr = test.expr
		}
		expr, err := Parse(test.expr)
		if err != nil {
			t.Error(err) // Ошибка анализа
			continue
		}
		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		if got != test.want {
			t.Errorf("%s.Eval() в %v = %q, требуется %q\n", test.expr, test.env, got, test.want)
		}

	}
}

func TestString(t *testing.T) {
	tests := []string{
		"x",
		"1.23",
		"-x",
		"x+y",
		"x*y",
		"x/y",
		"x-y",
		"pow(x,2)",
		"sin(x+y)",
		"sqrt(x*x+y*y)",
	}
	for _, test := range tests {
		expr1, err := Parse(test)
		if err != nil {
			t.Errorf("Parse %s: %v", test, err)
			continue
		}
		s := expr1.String()
		expr2, err := Parse(s)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", s, err)
			continue
		}
		if expr1.String() != expr2.String() {
			t.Errorf("не эквивалентны: %s != %s", expr1.String(), expr2.String())
		}
	}
}
