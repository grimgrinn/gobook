package bzip

import (
	"bytes"
	"sync"
	"testing"
)

func TestConcurrentWrite(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	var wg sync.WaitGroup
	n := 100
	data := []byte("hello world\n")

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.Write(data)
		}()
	}

	wg.Wait()
	w.Close()

	// Проверяем, что данные сжаты
	if buf.Len() == 0 {
		t.Errorf("ожидались сжатые данные, получено 0 байт:")
	}
}
