package main

import (
	"fmt"
	memo "gobook/9/practice9/9.3"
	"io/ioutil"
	"net/http"
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

func main() {
	m := memo.New(httpGetBody)

	// Без отмены
	v1, err := m.Get("https://go.dev", nil)
	if err != nil {
		fmt.Println("error :", err)
	}
	fmt.Println("Без отмены:", len(v1.([]byte)))

	// С отменой
	cancel := make(chan struct{})
	go func() {
		time.Sleep(1 * time.Millisecond)
		close(cancel)
	}()

	v2, err := m.Get("https://go.dev", cancel)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("C отменой: ", len(v2.([]byte)))
	}

}
