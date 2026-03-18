package main

import "fmt"

// nonempty возвращает срез, содержащий только непустые строки.
// Содержимое базоваго массива при работе функции изменяется.
func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}

func nonempty2(strings []string) []string {
	out := strings[:0] // Срез нулевой длины из исходного среза
	for _, s := range strings {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

func main() {
	data := []string{"one", "", "three"}
	fmt.Printf("%q\n", nonempty(data))
	fmt.Printf("%q\n", data)

	data = []string{"one", "", "three"}
	data = nonempty(data)
	fmt.Printf("%q\n", data)

	data = []string{"one", "", "three"}
	fmt.Printf("%q\n", nonempty2(data))
	fmt.Printf("%q\n", data)

}
