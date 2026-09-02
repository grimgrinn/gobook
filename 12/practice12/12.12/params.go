package main

import (
	"fmt"
	"gobook/12/params"
	"net/http"
)

type User struct {
	Name       string `http:"name,requred,min=2,max=50"`
	Email      string `http:"emai,email,required"`
	Age        int    `http:"age,min=18,max=120"`
	CreditCard string `http:"card, creaditCard"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := params.Unpack(r, &u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "User: %+v\n", u)
}

func main() {
	http.HandleFunc("/user", handler)
	http.ListenAndServe(":8080", nil)
}
