package main

import (
	"fmt"
	"gobook/12/params"
	"net/http"
)

type SearchData struct {
	Labels     []string `http:"labels"`
	MaxResults int      `http:"max"`
	Exact      bool     `http:"x"`
}

// search реализует окончание URL /search.
func search(resp http.ResponseWriter, req *http.Request) {
	var data SearchData

	data.MaxResults = 10 // Значение по умолчанию
	if err := params.Unpack(req, &data); err != nil {
		http.Error(resp, err.Error(), http.StatusBadRequest) // 400
		return
	}
	//... оствшаяся часть обработчика ...
	fmt.Fprintf(resp, "Поиск: %+v\n", data)
}

func main() {
	http.HandleFunc("/search", search)
	http.ListenAndServe(":8000", nil)
}
