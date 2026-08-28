package main

import (
	"fmt"
	"gobook/12/params"
	"net/http"
)

// search реализует окончание URL /search.
func search(resp http.ResponseWriter, req *http.Request) {
	var data struct {
		Labels     []string `http:"1"`
		MaxResults int      `http:"max"`
		Exact      bool     `http:"x"`
	}

	data.MaxResults = 10 // Значение по умолчанию
	if err := params.Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest) // 400
		return
	}
	//... оствшаяся часть обработчика ...
	fmt.Fprintf(resp, "Поиск: %+v\n", data)
}
