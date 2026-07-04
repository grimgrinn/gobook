// Пакет memo обеспечивает безопасное с точки зрения
// параллельности запоминание функции типа Func.
package memo

// Memo кеширует результаты вызовов Func.
type Memo struct {
	f     Func
	cache map[string]result
}

// Func является типом функции с запоминанием.
type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]result)}
}

// Небезопасно с точки зрения параллельности!
func (memo *Memo) Get(key string) (interface{}, error) {
	res, ok := memo.cache[key]
	if !ok {
		res.value, res.err = memo.f(key)
		memo.cache[key] = res
	}
	return res.value, res.err
}
