package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string        `json:"body"`
		UserId uuid.NullUUID `json:"user_id"`
	}
	type returnVals struct {
		ID        uuid.UUID     `json:"id"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
		Body      string        `json:"body"`
		UserID    uuid.NullUUID `json:"user_id"`
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

	attempt := struct {
		Body   string
		UserID uuid.NullUUID
	}{
		Body:   strings.Join(dividedChirp, " "),
		UserID: chirp.UserId,
	}
	newChirp, err := cfg.dbQueries.CreateChirpByUser(r.Context(), attempt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, returnVals{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	})

}
