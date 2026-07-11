package reader

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
)

func init() {
	fmt.Println("регистрация тар")
	RegisterFormat("tar", matchTar, newTar)
}

// matchTar проверяет TAR-магию (uStar)
func matchTar(r io.Reader) (bool, error) {
	// TAR не имеет фиксированной магии, но uStar имеет смещение 257: "ustart\x00"
	buf := make([]byte, 512)
	if _, err := io.ReadFull(r, buf); err != nil {
		return false, err
	}
	fmt.Printf("DEBUG TAR: первые 20 байт: %x\n", buf[:20])
	fmt.Printf("DEBUG TAR: на смещении 257: %x\n", buf[257:257+8])

	// Проверяем USTAR магию (GNU tar)
	if bytes.HasPrefix(buf[257:257+8], []byte("ustar\x00")) ||
		bytes.HasPrefix(buf[257:257+8], []byte("ustar \x00")) ||
		bytes.HasPrefix(buf[257:257+8], []byte("ustar  \x00")) {
		return true, nil
	}

	return false, nil
}

// tarReader обертка для tar.Reader
type tarReader struct {
	tr  *tar.Reader
	f   *tar.Header
	err error
}

func newTar(r io.ReaderAt, size int64) (ArchiveReader, error) {
	reader := io.NewSectionReader(r, 0, size)
	tr := tar.NewReader(reader)
	return &tarReader{tr: tr}, nil
}

func (t *tarReader) Next() bool {
	hdr, err := t.tr.Next()
	if err == io.EOF {
		return false
	}
	if err != nil {
		t.err = err
		return false
	}
	t.f = hdr
	return true
}

func (t *tarReader) File() File {
	if t.f == nil {
		return nil
	}
	return tarFile{t.f, t.tr}
}

func (t *tarReader) Err() error { return t.err }

type tarFile struct {
	hdr *tar.Header
	tr  *tar.Reader
}

func (t tarFile) Name() string { return t.hdr.Name }
func (t tarFile) Size() int64  { return t.hdr.Size }
func (t tarFile) Open() (io.ReadCloser, error) {
	return io.NopCloser(t.tr), nil
}
