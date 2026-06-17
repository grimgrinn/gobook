package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var verbose = flag.Bool("v", false, "показывать прогресс")

type result struct {
	root string
	size int64
	n    int64
}

type dirStats struct {
	files int64
	bytes int64
}

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	fileSizes := make(chan result)

	var wg sync.WaitGroup
	for _, root := range roots {
		wg.Add(1)
		go func(root string) {
			walkDir(root, root, fileSizes)
			wg.Done()
		}(root)
	}

	go func() {
		wg.Wait()
		close(fileSizes)
	}()

	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(1 * time.Second)
	}

	totals := make(map[string]dirStats)

	done := make(chan struct{})
	go func() {
		for res := range fileSizes {
			t := totals[res.root]
			t.files += res.n
			t.bytes += res.size
			totals[res.root] = t
		}
		done <- struct{}{}
	}()

loop:
	for {
		select {
		case <-tick:
			printTotals(roots, totals)
		case <-done:
			break loop
		}
	}
	printTotals(roots, totals)
}

func walkDir(root, dir string, fileSizes chan<- result) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkDir(root, subdir, fileSizes)
		} else {
			info, err := entry.Info()
			if err == nil {
				fileSizes <- result{root, info.Size(), 1}
			}
		}
	}
}

func dirents(dir string) []fs.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}

func printTotals(roots []string, totals map[string]dirStats) {
	for _, root := range roots {
		t := totals[root]
		fmt.Printf("%s: %d файлов, %.1f MB\n", root, t.files, float64(t.bytes)/1e6)
	}
}
