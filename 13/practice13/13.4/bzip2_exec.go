// Пакет bzip предоставляет writer, который использует сжатие bzip2
// через внешнюю утилиту bzip2.
package bzip

import (
	"io"
	"os/exec"
)

type writer struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	closed bool
}

// NewWriter возвращает writer для сжатых потоков.
// Использовать внешнюю утилиту bzip2.
func NewWriter(out io.Writer) io.WriteCloser {
	cmd := exec.Command("bzip2", "-c")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	w := &writer{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
	}

	go func() {
		io.Copy(out, stdout)
	}()

	return w
}

func (w *writer) Write(data []byte) (int, error) {
	if w.closed {
		panic("запись в закрытый поток")
	}
	return w.stdin.Write(data)
}

func (w *writer) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true

	if err := w.stdin.Close(); err != nil {
		return err
	}

	if err := w.cmd.Wait(); err != nil {
		return err
	}

	return nil
}
