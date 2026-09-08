//go:build cgo

// Пакет bzip предоставляет writer, который использует сжатие bzip2
package bzip

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS: -L/usr/lib -lbz2
#include <bzlib.h>
#include <stdlib.h>
int bz2compress(bz_stream *s, int action, char *in, unsigned *inlen, char *out, unsigned *outlen);
*/
import "C"
import (
	"io"
	"sync"
	"unsafe"
)

type writer struct {
	w      io.Writer
	stream *C.bz_stream
	outbuf *C.char
	outlen int
	mu     sync.Mutex
	closed bool
}

func NewWriter(out io.Writer) io.WriteCloser {
	const (
		blockSize  = 9
		verbosity  = 0
		workFactor = 30
	)
	bufSize := 64 * 1024
	w := &writer{
		w:      out,
		stream: new(C.bz_stream),
		outbuf: (*C.char)(C.malloc(C.size_t(bufSize))),
		outlen: bufSize,
	}

	C.BZ2_bzCompressInit(w.stream, blockSize, verbosity, workFactor)
	return w
}

func (w *writer) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		panic("запись в зыкрытый поток")
	}

	if w.stream == nil {
		panic("закрыт")
	}

	var total int

	for len(data) > 0 {
		inlen := C.uint(len(data))
		outlen := C.uint(w.outlen)

		// C.malloc для безопасной работый с памятью
		cdata := C.malloc(C.size_t(len(data)))
		if cdata == nil {
			return total, io.ErrShortWrite
		}
		defer C.free(cdata)

		copy((*[1 << 30]byte)(cdata)[:len(data)], data)

		C.bz2compress(w.stream, C.BZ_RUN,
			(*C.char)(cdata), &inlen,
			w.outbuf, &outlen)

		total += int(inlen)
		data = data[inlen:]

		if outlen > 0 {
			if _, err := w.w.Write(C.GoBytes(unsafe.Pointer(w.outbuf), C.int(outlen))); err != nil {
				return total, err
			}
		}
	}
	return total, nil
}

func (w *writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return nil // уже закрыт
	}
	w.closed = true

	if w.stream == nil {
		panic("закрыт")
	}

	for {
		inlen := C.uint(0)
		outlen := C.uint(w.outlen)

		var empty *C.char
		r := C.bz2compress(w.stream, C.BZ_FINISH, empty, &inlen,
			w.outbuf, &outlen)

		if outlen > 0 {
			if _, err := w.w.Write(C.GoBytes(unsafe.Pointer(w.outbuf), C.int(outlen))); err != nil {
				return err
			}
		}
		if r == C.BZ_STREAM_END {
			break
		}
	}

	C.BZ2_bzCompressEnd(w.stream)
	w.stream = nil

	if w.outbuf != nil {
		C.free(unsafe.Pointer(w.outbuf))
		w.outbuf = nil
	}

	return nil
}
