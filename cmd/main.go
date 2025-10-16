package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	port := flag.String("port", ":4000", "port the server runs on")

	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	srv := &http.Server{
		Addr:    *port,
		Handler: mux,
	}

	log.Printf("Listening on port %s", *port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
