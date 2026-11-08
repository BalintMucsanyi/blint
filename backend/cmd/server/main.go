package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := struct{ Title string }{Title: "Home"}
	render(w, "home", data)
}

func addSite(db *sql.DB, site Site) error {
	_, err := db.Exec("INSERT INTO sites (name, URL, Description) VALUES (?, ?, ?)", site.Name, site.URL, site.Description)
	return err
}

func addSiteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("Name")
		url := r.FormValue("URL")
		description := r.FormValue("Description")

		site := Site{Name: name, URL: url, Description: description}
		if err := addSite(db, site); err != nil {
			log.Println(err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/sites", http.StatusSeeOther)
	}
}

func sites(db *sql.DB) ([]Site, error) {
	rows, err := db.Query("SELECT id, name, URL, description FROM sites ORDER BY name DESC, id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Site
	for rows.Next() {
		var site Site
		if err := rows.Scan(&site.ID, &site.Name, &site.URL, &site.Description); err != nil {
			return nil, err
		}
		list = append(list, site)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func deleteSite(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM sites WHERE id = ?", id)
	return err
}

func deleteSiteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := deleteSite(db, id); err != nil {
			log.Println(err)
			http.Error(w, "error deleting site", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/sites", http.StatusSeeOther)
	}
}

func sitesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := sites(db)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		data := struct {
			Title string
			Sites []Site
		}{
			Title: "Sites",
			Sites: list,
		}
		render(w, "sites", data)
	}
}

type Site struct {
	ID          int
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
	rows, err := db.Query("SELECT id, date, name, duration, effort FROM activities ORDER BY date DESC, id DESC")
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

func deleteActivity(db *sql.DB, id int) error {
	if _, err := db.Exec("DELETE from activities WHERE id = ?", id); err != nil {
		return err
	}
	return nil
}

func deleteActivityHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := deleteActivity(db, id); err != nil {
			log.Println(err)
			http.Error(w, "error deleting activity", http.StatusInternalServerError)
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
