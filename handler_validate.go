package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := parameters{}
	err := decoder.Decode(&chirp)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	dividedChirp := strings.Split(chirp.Body, " ")
	for i, word := range dividedChirp {
		if ok := slices.Contains(BadWords, strings.ToLower(word)); ok {
			dividedChirp[i] = "****"
		}
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: strings.Join(dividedChirp, " "),
	})

}
