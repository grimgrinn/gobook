package bank2

var (
	sema    = make(chan struct{}, 1) // Бинарный семафор для
	balance int                      // защиты balace
)

func Deposit(amount int) {
	sema <- struct{}{} // Захват маркера
	balance = balance + amount
	<-sema // Освобождения маркера
}

func Balance() int {
	sema <- struct{}{} // Захват маркера
	b := balance
	<-sema // Освобождение маркера
	return b
}
