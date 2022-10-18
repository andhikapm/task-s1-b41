package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Main"))
	})

	route.HandleFunc("/alpha", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("alpha"))
	})

	route.HandleFunc("/beta", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("beta"))
	})

	fmt.Println("Server running on port 3000")
	http.ListenAndServe("localhost:3000", route)

}
