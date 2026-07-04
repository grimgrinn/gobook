package memo

type entry struct {
	res   result
	ready chan struct{} // Закрывается, когда res готов
}

// request представляет собой сообщение,
// требующее применения Func к key.
type request struct {
	key      string
	response chan<- result
	// Клиенту нужен только result
}

type Memo struct {
	requests chan request
}

type Func func(key string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// Это первый запрос данного ключа key.
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key) // Вызов f(key)
		}
		go e.deliver(req.response)
	}
}

func (e *entry) call(f Func, key string) {
	// Вычисление функции.
	e.res.value, e.res.err = f(key)
	// Оповещение о готовности
	close(e.ready)
}

func (e *entry) deliver(response chan<- result) {
	// Ожидание готовности.
	<-e.ready
	// Отправка результата клиенту.
	response <- e.res
}
