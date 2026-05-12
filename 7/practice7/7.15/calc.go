package main

import (
	"bufio"
	"fmt"
	"gobook/7/eval"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Введите выражение: ")
	if !scanner.Scan() {
		log.Fatal("Не удалось прочитать выражение")
	}

	exprStr := scanner.Text()
	if exprStr == "" {
		log.Fatal("Выражение не может быть пустым")
	}

	expr, err := eval.Parse(exprStr)
	if err != nil {
		log.Fatalf("Ошибка парсинга: %v\n", err)
	}

	// Проверяем выражение и ищем переменные
	vars := make(map[eval.Var]bool)
	if err := expr.Check(vars); err != nil {
		log.Fatalf("Ошибка проверки выражения: %v\n", err)
	}

	// Если нет переменных, не спрашиваем
	if len(vars) == 0 {
		result := expr.Eval(eval.Env{})
		fmt.Printf("Результат: %v\n", result)
		return
	}

	env := make(eval.Env)
	for v := range vars {
		fmt.Printf("Введите значение %s: ", v)
		if !scanner.Scan() {
			log.Fatalf("Не удалось прочитать значение для %s", v)
		}
		text := strings.TrimSpace(scanner.Text())
		val, err := strconv.ParseFloat(text, 64)
		if err != nil {
			log.Fatalf("Некорректное число для %s: %s", v, text)
		}
		env[v] = val
	}

	result := expr.Eval(env)
	fmt.Printf("Результат: %v\n", result)
}
