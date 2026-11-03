package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func pathHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "get path\n")
}
func sitesHandler(w http.ResponseWriter, _ *http.Request) {
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

func activitiesHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/layout.html", "templates/activities.html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	data := struct {
		Title      string
		Activities []Activity
	}{
		Title:      "Activities",
		Activities: activities(),
	}
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err)
	}
}

func activities() []Activity {
	return []Activity{
		{Name: "Bouldering", Date: "2026.10.28", Duration: "2h", Effort: "6"},
		{Name: "Street Workout", Date: "2026.11.07", Duration: "1h", Effort: "6"},
		{Name: "Bouldering", Date: "2026.11.04", Duration: "2h", Effort: "8"},
	}
}

type Activity struct {
	Date     string
	Name     string
	Duration string
	Effort   string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /path/", pathHandler)
	mux.HandleFunc("GET /sites", sitesHandler)
	// /{$} means "exactly / and nothing else"
	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("GET /activities", activitiesHandler)

	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatal(err)
	}
}
