package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Lando-Iraola/chirpy/internal/auth"
	"github.com/Lando-Iraola/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token", fmt.Errorf("Request with no token"))
		return
	}

	userID, err := auth.ValidateJWT(bearer, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", fmt.Errorf("Request with invalid token"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	userData := parameters{}
	err = decoder.Decode(&userData)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashed, err := auth.HashPassword(userData.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate secure password", err)
		return
	}

	updatedUser, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		HashedPassword: hashed,
		Email:          userData.Email,
		ID:             userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	})

}
