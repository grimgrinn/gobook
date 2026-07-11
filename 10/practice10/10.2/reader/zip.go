package reader

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

func init() {
	fmt.Println("регистрация zip")
	RegisterFormat("zip", matchZip, newZip)
}

// matchZip проверяет ZIP-магию (PK\x03\x04)
func matchZip(r io.Reader) (bool, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return false, err
	}
	fmt.Printf("DEBUG ZIP: %x\n", buf) // 504b0304
	return bytes.Equal(buf, []byte("PK\x03\x04")), nil
}

// zipReader обертка для zip.Reader
type zipReader struct {
	r   *zip.Reader
	idx int
	err error
}

func newZip(r io.ReaderAt, size int64) (ArchiveReader, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	return &zipReader{r: zr, idx: -1}, nil
}

func (z *zipReader) Next() bool {
	z.idx++
	return z.idx < len(z.r.File)
}

func (z *zipReader) File() File {
	if z.idx < 0 || z.idx >= len(z.r.File) {
		return nil
	}
	return zipFile{z.r.File[z.idx]}
}

func (z *zipReader) Err() error { return z.err }

type zipFile struct {
	f *zip.File
}

func (z zipFile) Name() string { return z.f.Name }
func (z zipFile) Size() int64  { return int64(z.f.UncompressedSize64) }
func (z zipFile) Open() (io.ReadCloser, error) {
	return z.f.Open()
}
