package main

import (
	"fmt"
	"gobook/7/eval"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tmpl *template.Template

func init() {
	// Загружаем щаблон при старте
	tmpl = template.Must(template.ParseFiles("templates/calc.html"))
}

func calcHandler(w http.ResponseWriter, r *http.Request) {
	// Данные для шаблона
	data := struct {
		Expr   string
		Vars   []string    // список имен переменных
		Result interface{} //результат (число или ошибка)
		Error  string
	}{}

	// Получаем выражение
	exprStr := r.FormValue("expr")
	data.Expr = exprStr

	if exprStr == "" {
		// ПОказываем пустую форму
		tmpl.Execute(w, data)
		return
	}

	// Парсим выражение
	expr, err := eval.Parse(exprStr)
	if err != nil {
		data.Error = fmt.Sprintf("Ошибка парсинга: %v", err)
		tmpl.Execute(w, data)
		return
	}

	// Находим переменные
	varsMap := make(map[eval.Var]bool)
	if err := expr.Check(varsMap); err != nil {
		data.Error = fmt.Sprintf("Ошибка проверки: %v", err)
		tmpl.Execute(w, data)
		return
	}

	// Преобразуем карту в срез для шаблона
	for v := range varsMap {
		data.Vars = append(data.Vars, string(v))
	}

	// Если это GET-запрос (первая загрузка страницы) - показываем форму с переменными
	if r.Method == http.MethodGet {
		tmpl.Execute(w, data)
		return
	}

	// POST-запрос: вычисляем выражение
	env := make(eval.Env)
	for _, v := range data.Vars {
		valStr := r.FormValue(v)
		if valStr == "" {
			data.Error = fmt.Sprintf("Переменная %s не заполнена", v)
			tmpl.Execute(w, data)
			return
		}
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			data.Error = fmt.Sprintf("Неверное число для %s:%s", v, valStr)
			tmpl.Execute(w, data)
			return
		}
		env[eval.Var(v)] = val
	}

	result := expr.Eval(env)
	data.Result = result
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", calcHandler)
	log.Println("Сервер запущен на http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
