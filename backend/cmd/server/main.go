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

func routes(store Store) *http.ServeMux {
	mux := http.NewServeMux()

	// /{$} means "exactly / and nothing else"
	mux.HandleFunc("GET /{$}", homeHandler)

	mux.HandleFunc("GET /sites", sitesHandler(store.db))
	mux.HandleFunc("POST /sites", addSiteHandler(store.db))
	mux.HandleFunc("POST /sites/{id}/delete", deleteSiteHandler(store.db))

	mux.HandleFunc("GET /activities", activitiesHandler(&store))
	mux.HandleFunc("POST /activities", addActivityHandler(&store))
	mux.HandleFunc("POST /activities/{id}/delete", deleteActivityHandler(&store))

	return mux
}

func initDB(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS activities (
    id       INTEGER PRIMARY KEY,
    date     TEXT    NOT NULL,
    name     TEXT    NOT NULL,
    duration INTEGER NOT NULL,
    effort   INTEGER NOT NULL CHECK (effort BETWEEN 1 AND 10)
);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS sites (
	id			INTEGER PRIMARY KEY,
	name		TEXT NOT NULL,
	URL			TEXT NOT NULL,
	Description TEXT NOT NULL
);`)
	if err != nil {
		return err
	}
	return nil
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

	if err := initDB(db); err != nil {
		log.Fatal(err)
	}
	store := NewStore(db)

	mux := routes(*store)

	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatal(err)
	}
}
