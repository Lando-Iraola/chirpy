package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Chirp struct {
	Body string `json:"body"`
}

func handlerChirp(w http.ResponseWriter, r *http.Request) {
	type Valid struct {
		Valid bool `json:"valid"`
	}
	type Error struct {
		Error string `json:"error"`
	}
	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)

	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		erBody := Error{
			Error: fmt.Sprintf("Error decoding parameters %s", err),
		}
		data, err := json.Marshal(erBody)
		if err != nil {
			log.Printf("Error marshaling data!")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	if len(chirp.Body) > 140 {
		erBody := Error{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(erBody)
		if err != nil {
			log.Printf("Error marshaling data!")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	data, err := json.Marshal(Valid{Valid: true})
	if err != nil {
		log.Printf("Error marshaling data!")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}
