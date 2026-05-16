package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Lando-Iraola/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
		return
	}

	if key != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	const upgradeEvent = "user.upgraded"

	decoder := json.NewDecoder(r.Body)
	event := parameters{}
	err = decoder.Decode(&event)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse request", fmt.Errorf("Failed to parse request"))
		return
	}

	if event.Event != upgradeEvent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.dbQueries.UpgradeUserToChirpyRed(r.Context(), event.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to parse request", fmt.Errorf("Failed to parse request"))
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
