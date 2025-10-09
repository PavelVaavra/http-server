package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/PavelVaavra/http-server/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createChirps(w http.ResponseWriter, r *http.Request) {
	const validLength = 140

	type chirpParams struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := chirpParams{}
	err := decoder.Decode(&params)
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
		UserID: params.UserId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
		return
	}

	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
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
