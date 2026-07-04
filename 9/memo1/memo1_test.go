package memo_test

import (
	"fmt"
	memo "gobook/9/memo1"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func TestMemo(t *testing.T) {
	urls := []string{
		"https://go.dev",
		"https://godoc.org",
		"https://play.golang.org",
	}

	m := memo.New(httpGetBody)

	for _, url := range urls {
		start := time.Now()
		value, err := m.Get(url)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%s, %s, %d байтов\n", url, time.Since(start), len(value.([]byte)))
	}

	for _, url := range urls {
		start := time.Now()
		value, err := m.Get(url)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%s, %s, %d байтов\n", url, time.Since(start), len(value.([]byte)))
	}
}

func TestConcurrentMemo(t *testing.T) {
	urls := []string{
		"https://go.dev",
		"https://godoc.org",
		"https://play.golang.org",
	}

	m := memo.New(httpGetBody)

	var n sync.WaitGroup
	for _, url := range urls {
		n.Add(1)
		go func(u string) {
			start := time.Now()
			value, err := m.Get(u)
			if err != nil {
				t.Error(err)
			}
			fmt.Printf("%s, %s, %d байтов\n", u, time.Since(start), len(value.([]byte)))
			n.Done()
		}(url)
	}
	n.Wait()
}
