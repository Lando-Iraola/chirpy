package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Lando-Iraola/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string        `json:"email"`
		Password         string        `json:"password"`
		ExpiresInSeconds time.Duration `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	userCreds := parameters{}
	err := decoder.Decode(&userCreds)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), userCreds.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	isSame, err := auth.CheckPasswordHash(userCreds.Password, user.HashedPassword)
	if err != nil || !isSame {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	duration := userCreds.ExpiresInSeconds
	if userCreds.ExpiresInSeconds == 0 || userCreds.ExpiresInSeconds > time.Duration(1*time.Hour) {
		duration = time.Duration(1 * time.Hour)
	}

	t, err := auth.MakeJWT(user.ID, cfg.secret, duration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     t,
	})
}
