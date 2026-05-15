package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't query", err)
		return
	}

	manyChirps := make([]Chirp, 0)
	for _, chirp := range chirps {
		manyChirps = append(manyChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, manyChirps)
}

func (cfg *apiConfig) handlerChirpsRetrieveOne(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("chirpID")

	if len(rawId) != 36 {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", fmt.Errorf("No chirp found by id: %s", rawId))
		return
	}
	id := uuid.MustParse(rawId)

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't query", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
