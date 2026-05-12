package main

import (
	"log"
	"net/http"
)

func rootHandler() http.Handler {

	handler := http.FileServer(http.Dir("."))
	return handler
}

func main() {
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/", rootHandler())

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
