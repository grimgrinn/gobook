package memo

import "fmt"

type entry struct {
	res   result
	ready chan struct{}
}

type request struct {
	key      string
	done     <-chan struct{}
	response chan<- result
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

func (memo *Memo) Get(key string, done <-chan struct{}) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, done, response}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests {
		// Проверяем не отменен ли запрос
		if req.done != nil && cancelled(req.done) {
			req.response <- result{err: fmt.Errorf("request cancelled")}
			continue
		}

		e := cache[req.key]
		if e == nil {
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key, req.done)
		}
		go e.deliver(req.response, req.done)
	}
}

func (e *entry) call(f Func, key string, done <-chan struct{}) {
	// Выполняем функциюб но передаем ей канал отмены
	e.res.value, e.res.err = f(key)

	// Если запрос был отменен, не закрываем ready (кеширование отменяетсяФ)
	if done != nil && cancelled(done) {
		// Оставляем entry пустым (не сохрвняем результат)
		return
	}

	close(e.ready)
}

func (e *entry) deliver(response chan<- result, done <-chan struct{}) {
	if done == nil {
		// Отмена не поддерживается - прсто ждем готовности
		<-e.ready
		response <- e.res
		return
	}

	select {
	case <-e.ready:
		// Результат готов - отправляем
		response <- e.res
	case <-done:
		// Запрос отменен - отправляем ошибку
		response <- result{err: fmt.Errorf("request cancelled")}
	}
}

// cancelled проверяет, не отменен ли запрос
func cancelled(done <-chan struct{}) bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}
