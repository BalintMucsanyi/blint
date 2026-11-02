package main

import (
	"fmt"
	"log"
	"net/http"
)

func pathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "get path\n")
}
func helloHandler(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello\n")
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home\n")
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /path/blint", pathHandler)
	mux.HandleFunc("GET /hello", helloHandler)
	// /{$} means "exactly / and nothing else"
	mux.HandleFunc("GET /{$}", homeHandler)

	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatal(err)
	}
}
