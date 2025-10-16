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
