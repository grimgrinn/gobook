package main

import "fmt"

type tree struct {
	value       int
	left, right *tree
}

func (t *tree) String() string {
	var values []int
	var inorder func(*tree)
	inorder = func(node *tree) {
		if node == nil {
			return
		}
		inorder(node.left)
		values = append(values, node.value)
		inorder(node.right)
	}
	inorder(t)
	return fmt.Sprint(values)
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
	root := &tree{value: 5}
	root = add(root, 2)
	root = add(root, 7)
	root = add(root, 1)
	root = add(root, 3)
	fmt.Println(root) // [1 2 3 5 7]
}
