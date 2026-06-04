package main

import (
	"log"
	"os"
	"sync"

	"gopl.io/ch8/thumbnail"
)

// makeThumbnails создает эскизы для файлов из списка.
func makeThumbnails(filenames []string) {
	for _, f := range filenames {
		if _, err := thumbnail.ImageFile(f); err != nil {
			log.Println(err)
		}
	}
}

// Неправильно!
func makeThumbnails2(filenames []string) {
	for _, f := range filenames {
		go thumbnail.ImageFile(f) // Примечание: игнорируем ошибки
	}
}

// makeThumbnails3 параллельно делает эскизы определенных файлов.
func makeThumbnails3(filenames []string) {
	ch := make(chan struct{})
	for _, f := range filenames {
		go func(f string) {
			thumbnail.ImageFile(f) // Примечание: игнорируем ошибки
			ch <- struct{}{}
		}(f)
	}
	// Ожидание завершения go-подпрограммы
	for range filenames {
		<-ch
	}
}

// makeTubmnails4 параллельно делает эскизы определенных файлов.
// Возвращает ошиьку при сбое на любом шаге.
// func makeThumbnails4(filenames []string) {
// 	errors := make(chan error)

// 	for _, f := range filenames {
// 		go func(f string) {
// 			_, err := thumbnail.ImageFile(f)
// 			errors <- err
// 		}(f)
// 	}

// 	for range filenames {
// 		if err := <-errors; err != nil {
// 			return err // Примечание: неверно: утечка go-подпрограмм!
// 		}
// 	}
// 	return nil
// }

// makeThumbnails5 параллельно делает эскизы определенных файлов.
// Возвращает имена сгенерированных файлов в произвольном порядке
// или ошибку при сбое на любом шаге.
func makeThumbnails5(filenames []string) (thumbfiles []string, err error) {
	type item struct {
		thumbfile string
		err       error
	}

	ch := make(chan item, len(filenames))
	for _, f := range filenames {
		go func(f string) {
			var it item
			it.thumbfile, it.err = thumbnail.ImageFile(f)
		}(f)
	}

	for range filenames {
		it := <-ch
		if it.err != nil {
			return nil, it.err
		}
		thumbfiles = append(thumbfiles, it.thumbfile)
	}

	return thumbfiles, nil
}

// makeThumbnails6 параллельно делает эскизы определенных файлов.
// Возвращает общее количество байтов в созданных файлах.
func makeThumbnails6(filenames <-chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup // количество работающих go-подпрограмм
	for f := range filenames {
		wg.Add(1)
		// Рабочая go-программа
		go func(f string) {
			defer wg.Done()
			thumb, err := thumbnail.ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}
			info, _ := os.Stat(thumb) // Игнорируем ошибки
			sizes <- info.Size()
		}(f)
	}

	// Ожидание счетчика
	go func() {
		wg.Wait()
		close(sizes)
	}()

	var total int64
	for size := range sizes {
		total += size
	}
	return total
}
