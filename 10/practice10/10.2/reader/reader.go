// Пакет reader предоставялет обобщенный интерфейс для чтения архивов.
package reader

import (
	"bytes"
	"fmt"
	"io"
	"log"
)

// File предствляет один файл в архиве
type File interface {
	Name() string
	Size() int64
	Open() (io.ReadCloser, error)
}

// ArchiveReader представляет архив с возможностью итерации по файлам.
type ArchiveReader interface {
	Next() bool
	File() File
	Err() error
}

// Format описывает формат архива.
type Format struct {
	Name  string
	Match func(r io.Reader) (bool, error)
	New   func(r io.ReaderAt, size int64) (ArchiveReader, error)
}

var formats []Format

func RegisterFormat(name string, match func(io.Reader) (bool, error), new func(io.ReaderAt, int64) (ArchiveReader, error)) {
	formats = append(formats, Format{Name: name, Match: match, New: new})
}

// Open открывает архив и возвращает ArchiveReader.
func Open(r io.ReaderAt, size int64) (ArchiveReader, string, error) {
	// Создаем io.Reader для чтения заголовка

	fmt.Printf("DEBUG: formrats ocunt = %d\n", len(formats))
	header := make([]byte, 512)
	if _, err := io.ReadFull(io.NewSectionReader(r, 0, size), header); err != nil {
		return nil, "", err
	}

	for _, fmt := range formats {
		log.Printf("Debug: проверяю формат %s\n", fmt.Name)

		reader := bytes.NewReader(header)

		ok, err := fmt.Match(reader)
		if err != nil {
			log.Printf("DEBUG: ошибка в match: %v\n", err)
			continue
		}
		if ok {
			log.Printf("Debug: формат %s определен\n", fmt.Name)
			ar, err := fmt.New(r, size)
			if err != nil {
				return nil, "", err
			}
			return ar, fmt.Name, nil
		}
	}
	return nil, "", fmt.Errorf("неизвестный формат архива")
}
