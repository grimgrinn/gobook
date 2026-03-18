package main

import "fmt"

func strUniq(str []string) []string {
	if len(str) <= 1 {
		return str
	}

	write := 1 // следующая позиция для записи уникального элемента

	for read := 1; read < len(str); read++ {
		if str[read] != str[write-1] {
			str[write] = str[read]
			write++
		}
	}
	return str[:write]
}

func main() {

	s := []string{"test", "sdfs", "sdfs", "sdfsdf", "234", "123", "123", "234", "2", "2", "23432", "", "2344"}
	fmt.Println(s)

	s = strUniq(s)

	fmt.Println(s)

}
