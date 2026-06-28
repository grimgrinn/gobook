// Пакет bank предоставляет безопасный с точки зрения
// параллельности банк с одним счетом.
package bank

var deposits = make(chan int) // Отправление вклада
var balances = make(chan int) // Получение баланса

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func teller() {
	var balance int // balance ограничен go-подпрограммой teller
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // Запуск управляющей go-подпрограммы
}
