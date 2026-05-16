package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Lando-Iraola/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token given", fmt.Errorf("No token given"))
		return
	}

	token, err := cfg.dbQueries.FindToken(r.Context(), bearer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token found", fmt.Errorf("No token found"))
		return
	}
	if token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token expired", fmt.Errorf("Token expired"))
		return
	}

	if time.Now().After(token.ExpiresAt) {
		cfg.dbQueries.RevokeRefreshToken(r.Context(), token.Token)
		respondWithError(w, http.StatusUnauthorized, "Token expired", fmt.Errorf("Token expired"))
		return
	}

	newToken, err := auth.MakeJWT(token.UserID, cfg.secret, time.Duration(time.Hour))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create token", fmt.Errorf("Failed to create token"))
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{Token: newToken})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token given", fmt.Errorf("No token given"))
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), bearer)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token", fmt.Errorf("Error revoking token"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
