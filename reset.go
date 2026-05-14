package main

import (
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if platform := os.Getenv("PLATFORM"); platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(http.StatusText(http.StatusForbidden)))
		return
	}

	err := cfg.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error when trying to clear all users: %s", err)
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
