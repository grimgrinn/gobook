package main

import (
	"fmt"
)

func join(sep string, strs ...string) string {
	switch len(strs) {
	case 0:
		return ""
	case 1:
		return strs[0]
	}

	result := ""

	for i, str := range strs {
		if i > 0 {
			result += sep
		}
		result += str

	}

	return result
}

func main() {

	fmt.Println(join(", ", "jopa", "siska", "pizda", "piska"))

}
