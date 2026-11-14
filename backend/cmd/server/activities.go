package main

import (
	"log"
	"net/http"
	"strconv"
)

func activitiesHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := store.activities()
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

func (s *Store) activities() ([]Activity, error) {
	rows, err := s.db.Query("SELECT id, date, name, duration, effort FROM activities ORDER BY date DESC, id DESC")
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
func (s *Store) AddActivity(a Activity) error {
	_, err := s.db.Exec("INSERT INTO activities (date, name, duration, effort) VALUES (?, ?, ?, ?)", a.Date, a.Name, a.Duration, a.Effort)
	return err
}

func addActivityHandler(store *Store) http.HandlerFunc {
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
		if err := store.AddActivity(a); err != nil {
			log.Println(err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/activities", http.StatusSeeOther)
	}
}

func (s *Store) DeleteActivity(id int) error {
	if _, err := s.db.Exec("DELETE from activities WHERE id = ?", id); err != nil {
		return err
	}
	return nil
}

func deleteActivityHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := store.DeleteActivity(id); err != nil {
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
