package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	port := flag.String("port", ":4000", "port the server runs on")

	srv := &http.Server{
		Addr: *port,
	}

	log.Printf("Listening on port %s", *port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
