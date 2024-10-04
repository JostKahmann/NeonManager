package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {

	handler1 := func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)
	}

	http.HandleFunc("/", handler1)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
