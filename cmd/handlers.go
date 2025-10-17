package main

import (
	"html/template"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/home.tmpl.html",
		"./ui/html/ascii.tmpl.html",
	}

	// read each tempalte file and parse it
	tempSet, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tempSet.ExecuteTemplate(w, "home", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func blog(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/blog" {
		http.NotFound(w, r)
		return
	}

	// read each tempalte file and parse it
	tempSet, err := template.ParseFiles("./ui/html/blog.tmpl.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tempSet.ExecuteTemplate(w, "blog", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
