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
	for _, p := range sites() {
		fmt.Fprintf(w, "%s - %s (%s)\n", p.Name, p.Description, p.URL)
	}
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home\n")
}

type FavouriteSite struct {
	Name        string
	Description string
	URL         string
}

func sites() []FavouriteSite {
	return []FavouriteSite{
		{Name: "boardgamegeek", Description: "A place where you can browse great boardgames", URL: "https://boardgamegeek.com/"},
		{Name: "krillion.io", Description: "3 min long daily game", URL: "https://krillion.io"},
		{Name: "localhost", Description: "Different for everyone", URL: "http://127.0.0.1"},
	}
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
