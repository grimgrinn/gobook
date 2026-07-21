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

func TestMin(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want float64
	}{
		{"min(1,2,3)", nil, 1},
		{"min(5,4,3,2,1)", nil, 1},
		{"min(x,y,z)", Env{"x": 10, "y": 20, "z": 5}, 5},
		{"min(10, x, 30)", Env{"x": 25}, 10},
		{"min(x,y) + min(a,b)", Env{"x": 100, "y": 200, "a": 50, "b": 150}, 150}, // 100 + 50
	}

	for _, test := range tests {
		expr, err := Parse(test.expr)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", test.expr, err)
			continue
		}

		got := expr.Eval(test.env)
		if got != test.want {
			t.Errorf("%s.Eval() = %v, want %v", test.expr, got, test.want)
		}
	}
}

func TestMinString(t *testing.T) {
	expr, _ := Parse("min(1, x, 3)")
	want := "min(1, x, 3)"
	if expr.String() != want {
		t.Errorf("String() = %q, want %q", expr.String(), want)
	}
}

func TestMinError(t *testing.T) {
	_, err := Parse("min()")
	if err == nil {
		t.Error("Parse(\"min()\") должно вернуть ошибку, но вернуло nil")
	}

}

func TestCoverage(t *testing.T) {
	var tests = []struct {
		input string
		env   Env
		want  string // Ожидаемая ошибка от Parse/Check
		// или результат Eval
	}{
		{"x % 2", nil, "unexpected '%'"},
		{"!true", nil, "unexpected '!'"},
		{"log(10)", nil, `неизвестная функция "log"`},
		{"sqrt(1, 2)", nil, "вызов sqrt имеет 2 вместо 1 аргументов"},
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
	}
	for _, test := range tests {
		expr, err := Parse(test.input)
		if err == nil {
			err = expr.Check(map[Var]bool{})
		}
		if err != nil {
			if err.Error() != test.want {
				t.Errorf("%s: получено %q, требуется %q", test.input, err, test.want)
			}
			continue
		}
		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		if got != test.want {
			t.Errorf("%s: %v => %s, want %s", test.input, test.env, got, test.want)
		}
	}
}
