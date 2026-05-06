package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type database struct {
	sync.Mutex
	items map[string]dollars
}
type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db.items {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db.items[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	if item == "" || price == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "item and price required")
		return
	}

	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid price: %s", price)
		return
	}

	db.Lock()
	defer db.Unlock()
	db.items[item] = dollars(p)
	fmt.Fprintf(w, "updated: %s = %s\n", item, dollars(p))
}

func (db *database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if item == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "item required")
		return
	}

	db.Lock()
	defer db.Unlock()
	delete(db.items, item)
	fmt.Fprintf(w, "deleted: %s\n", item)
}

func main() {
	db := &database{
		items: map[string]dollars{"shoes": 50, "socks": 5},
	}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}
