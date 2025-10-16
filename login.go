package main

import (
	"encoding/json"
	"net/http"

	"github.com/PavelVaavra/http-server/internal/auth"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type userParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := userParams{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode JSON", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJson(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
