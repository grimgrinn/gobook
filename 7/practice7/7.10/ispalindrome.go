package main

import (
	"fmt"
	"sort"
)

func IsPalindrome(s sort.Interface) bool {
	n := s.Len()
	for i := 0; i < n/2; i++ {
		j := n - 1 - i
		if s.Less(i, j) || s.Less(j, i) {
			return false
		}
	}
	return true
}

func main() {
	data := sort.IntSlice{1, 2, 3, 2, 1}
	fmt.Println(IsPalindrome(data)) // true
	data2 := sort.IntSlice{1, 2, 4, 5, 4}
	fmt.Println(IsPalindrome(data2)) // false
	words := sort.StringSlice{"a", "b", "c", "b", "a"}
	fmt.Println(IsPalindrome(words)) // true

}
