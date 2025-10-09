package main

import "net/http"

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not list chirps.", err)
		return
	}

	payload := []Chirp{}

	for _, chirp := range chirps {
		payload = append(payload, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}

	respondWithJson(w, http.StatusOK, payload)
}
