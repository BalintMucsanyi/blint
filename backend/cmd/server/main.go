package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func sitesHandler(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Title string
		Sites []FavouriteSite
	}{
		Title: "Sites",
		Sites: sites(),
	}
	render(w, "sites", data)
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := struct{ Title string }{Title: "Home"}
	render(w, "home", data)
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
	data := struct {
		Title      string
		Activities []Activity
	}{
		Title:      "Activities",
		Activities: activities(),
	}
	render(w, "activities", data)
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
	db, err := sql.Open("sqlite", "blint.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS activities (
    id       INTEGER PRIMARY KEY,
    date     TEXT    NOT NULL,
    name     TEXT    NOT NULL,
    duration INTEGER NOT NULL,
    effort   INTEGER NOT NULL CHECK (effort BETWEEN 1 AND 10)
);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(
		`INSERT INTO activities (date, name, duration,
		effort) VALUES (?, ?, ?, ?)`,
		"2026-10-28", "Bouldering", "120", "6",
	)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /sites", sitesHandler)
	// /{$} means "exactly / and nothing else"
	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("GET /activities", activitiesHandler)

	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatal(err)
	}
}
