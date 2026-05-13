package main

import (
	"fmt"
	"gobook/7/eval"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func calcHandler(w http.ResponseWriter, r *http.Request) {
	exprStr := r.FormValue("expr")
	if exprStr == "" {
		render(w, "", nil, nil)
		return
	}

	expr, err := eval.Parse(exprStr)
	if err != nil {
		render(w, exprStr, nil, fmt.Errorf("ошибка парсинга: %v", err))
		return
	}

	vars := make(map[eval.Var]bool)
	if err := expr.Check(vars); err != nil {
		render(w, exprStr, nil, fmt.Errorf("ошибка проверки: %v", err))
		return
	}

	env := make(eval.Env)
	for v := range vars {
		val, err := strconv.ParseFloat(r.FormValue(string(v)), 64)
		if err != nil {
			render(w, exprStr, vars, fmt.Errorf("неверное значение для %s", v))
			return
		}
		env[v] = val
	}

	result := expr.Eval(env)
	render(w, exprStr, vars, result)
}

func render(w http.ResponseWriter, expr string, vars map[eval.Var]bool, data interface{}) {
	tmpl := template.Must(template.ParseFiles("templates/calc.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Expr":   expr,
		"Vars":   vars,
		"Result": data,
	})

}

func main() {
	http.HandleFunc("/", calcHandler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
