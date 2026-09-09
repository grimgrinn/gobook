package bzip

import (
	"bytes"
	"os"
	"os/exec"
	"sync"
	"testing"
)

func TestBzipExec(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	data := []byte("hello worold\n")
	if _, err := w.Write(data); err != nil {
		t.Fatal(err)
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	if buf.Len() == 0 {
		t.Error("ожидались сжатые данные, получено 0 байт")
	}

	if err := testBunzip2(buf.Bytes()); err != nil {
		t.Error(err)
	}
}

func testBunzip2(data []byte) error {
	f, err := os.CreateTemp("", "text.bz2")
	if err != nil {
		return err
	}

	defer os.Remove(f.Name())

	if _, err := f.Write(data); err != nil {
		return err
	}
	f.Close()

	cmd := exec.Command("bunzip2", "-c", f.Name())
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	if string(out) != "hello world\n" {
		return err
	}
	return nil
}

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

	if buf.Len() == 0 {
		t.Errorf("ожидались сжатые данные, получен 0 байт")
	}
}
