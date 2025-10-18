package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/PavelVaavra/http-server/internal/auth"
	"github.com/PavelVaavra/http-server/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirps(w http.ResponseWriter, r *http.Request) {
	const validLength = 140

	type chirpParams struct {
		Body string `json:"body"`
	}

	jwtToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User isn't logged in, no JWT token available", err)
		return
	}
	userId, err := auth.ValidateJWT(jwtToken, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User isn't logged in, JWT isn't validated", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := chirpParams{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode JSON", err)
		return
	}

	if len(params.Body) > validLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedBody := replaceProfaneWords(params.Body)

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
		return
	}

	respondWithJson(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func replaceProfaneWords(s string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(s, " ")

	for i, word := range words {
		loweredWord := strings.ToLower(word)
		for _, profaneWord := range profaneWords {
			if loweredWord == profaneWord {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}
