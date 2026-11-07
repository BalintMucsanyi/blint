package main

import (
	"html/template"
	"log"
	"net/http"
)

func render(w http.ResponseWriter, page string, data any) {
	t, err := template.ParseFiles("templates/layout.html", "templates/"+page+".html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err)
	}
}
