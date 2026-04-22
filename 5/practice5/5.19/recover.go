package main

import "fmt"

func recov() (result int) {
	defer func() {
		if v := recover(); v != nil {
			result = v.(int)
		}
	}()

	panic(100)
}

func main() {
	fmt.Println(recov()) // 100
}
