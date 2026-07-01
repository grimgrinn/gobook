package bank

var deposits = make(chan int)
var balances = make(chan int)
var withdrawals = make(chan withdrawRequest)

type withdrawRequest struct {
	amount int
	result chan bool
}

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }

func Withdraw(amount int) bool {
	result := make(chan bool)
	withdrawals <- withdrawRequest{amount, result}
	return <-result
}

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case req := <-withdrawals:
			if balance >= req.amount {
				balance -= req.amount
				req.result <- true
			} else {
				req.result <- false
			}
		}
	}
}

func init() {
	go teller()
}
