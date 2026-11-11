package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

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
