package popcount

import "sync"

var (
	pc   [256]byte
	once sync.Once
)

// initTable инициализирует таблицу поиска
func initTable() {
	for i := range pc {
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount возвращает количество утсановленных битов в x
func PopCount(x uint64) int {
	once.Do(initTable) // инициализация при первом вызове
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}
