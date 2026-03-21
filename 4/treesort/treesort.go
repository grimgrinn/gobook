package main

import "fmt"

type tree struct {
	value       int
	left, right *tree
}

// Sort сортирует значение "на месте".
func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

// appendValues добавляет элементы t к values в требуемом
// порядке и возвращает результирующий срез.
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		// Эквивалентно возврату &tree{value: value}
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

func main() {
	values := []int{5, 2, 7, 1, 3, 9, 4}
	fmt.Println("До:", values)
	Sort(values)
	fmt.Println("После:", values)

	empty := []int{}
	Sort(empty)
	fmt.Println("Пустой:", empty)

	single := []int{42}
	Sort(single)
	fmt.Println("Один:", single)

	sorted := []int{1, 2, 3, 4, 5}
	Sort(sorted)
	fmt.Println("Уже отсортированный:", sorted)

	reversed := []int{9, 8, 7, 6, 5}
	Sort(reversed)
	fmt.Println("Обратный:", reversed)

	duplicates := []int{3, 1, 4, 1, 5, 9, 2, 6, 5}
	Sort(duplicates)
	fmt.Println("С дубликатами:", duplicates)
}
