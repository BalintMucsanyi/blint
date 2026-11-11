package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := struct{ Title string }{Title: "Home"}
	render(w, "home", data)
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

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS sites (
	id			integer PRIMARY KEY,
	name		TEXT NOT NULL,
	URL			TEXT NOT NULL,
	Description TEXT NOT NULL
);`)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /sites", sitesHandler(db))
	mux.HandleFunc("POST /sites", addSiteHandler(db))
	mux.HandleFunc("POST /sites/{id}/delete", deleteSiteHandler(db))

	// /{$} means "exactly / and nothing else"
	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("GET /activities", activitiesHandler(db))
	mux.HandleFunc("POST /activities", addActivityHandler(db))
	mux.HandleFunc("POST /activities/{id}/delete", deleteActivityHandler(db))

	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatal(err)
	}
}
