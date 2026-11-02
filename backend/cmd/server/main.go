package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func pathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "get path\n")
}
func sitesHandler(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("templates/layout.html", "templates/sites.html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	data := struct {
		Title string
		Sites []FavouriteSite
	}{
		Title: "Sites",
		Sites: sites(),
	}

	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err)
	}
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/layout.html", "templates/home.html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	data := struct{ Title string }{Title: "Home"}

	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err)
	}
}

func sites() []FavouriteSite {
	return []FavouriteSite{
		{Name: "boardgamegeek", Description: "A place where you can browse great boardgames", URL: "https://boardgamegeek.com/"},
		{Name: "krillion.io", Description: "3 min long daily game", URL: "https://krillion.io"},
		{Name: "localhost", Description: "Different for everyone", URL: "http://127.0.0.1"},
	}
}

type FavouriteSite struct {
	Name        string
	Description string
	URL         string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /path/", pathHandler)
	mux.HandleFunc("GET /sites", sitesHandler)
	// /{$} means "exactly / and nothing else"
	mux.HandleFunc("GET /{$}", homeHandler)

	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatal(err)
	}
}
