package memo

import "sync"

type Memo struct {
	f     Func
	mu    sync.Mutex // Защита cache
	cache map[string]result
}

type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// Метод Get безопасен с точки зрения параллельности.
func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	res, ok := memo.cache[key]
	memo.mu.Unlock()
	if !ok {
		res.value, res.err = memo.f(key)

		// Между этими двумя критическими разделами
		// несколько go-полпрограмм могут вычислять
		// f(key) и обновлять отображеиие.
		memo.mu.Lock()
		memo.cache[key] = res
		memo.mu.Unlock()
	}
	return res.value, res.err
}
