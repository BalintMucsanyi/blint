package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

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

func activitiesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := activities(db)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		data := struct {
			Title      string
			Activities []Activity
		}{"Activities", list}
		render(w, "activities", data)
	}
}

func activities(db *sql.DB) ([]Activity, error) {
	rows, err := db.Query("SELECT id, date, name, duration, effort FROM activities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Activity

	for rows.Next() {
		var a Activity
		if err := rows.Scan(&a.ID, &a.Date, &a.Name, &a.Duration, &a.Effort); err != nil {
		}
		list = append(list, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}
func addActivity(db *sql.DB, a Activity) error {
	_, err := db.Exec("INSERT INTO activities (date, name, duration, effort) VALUES (?, ?, ?, ?)", a.Date, a.Name, a.Duration, a.Effort)
	return err
}

func addActivityHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		date := r.FormValue("date")
		name := r.FormValue("name")

		duration, err := strconv.Atoi(r.FormValue("duration"))
		if err != nil {
			http.Error(w, "duration must be a number", http.StatusBadRequest)
			return
		}
		effort, err := strconv.Atoi(r.FormValue("effort"))
		if err != nil {
			http.Error(w, "effort must be a number", http.StatusBadRequest)
			return
		}

		a := Activity{Date: date, Name: name, Duration: duration, Effort: effort}
		if err := addActivity(db, a); err != nil {
			log.Println(err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/activities", http.StatusSeeOther)
	}
}

type Activity struct {
	ID       int
	Date     string
	Name     string
	Duration int
	Effort   int
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /sites", sitesHandler)
	// /{$} means "exactly / and nothing else"
	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("GET /activities", activitiesHandler(db))
	mux.HandleFunc("POST /activities", addActivityHandler(db))

	if err := http.ListenAndServe(":8090", mux); err != nil {
		log.Fatal(err)
	}
}
