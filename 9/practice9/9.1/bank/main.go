package main

import (
	"fmt"
	bank "gobook/9/practice9/9.1"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		bank.Deposit(200)
		fmt.Println("Баланс после повышения:", bank.Balance())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ok := bank.Withdraw(50)
		fmt.Println("Снятие 50:", ok)
		fmt.Println("Баланс после снятия:", bank.Balance())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ok := bank.Withdraw(200)
		fmt.Println("Снятие 200:", ok)
		fmt.Println("Баланс после попытки снять 200:", bank.Balance())
	}()

	wg.Wait()
}
