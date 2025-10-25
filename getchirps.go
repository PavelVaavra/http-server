package main

import (
	"net/http"
	"sort"

	"github.com/PavelVaavra/http-server/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	id, err := uuid.Parse(s)
	if err != nil && s != "" {
		respondWithError(w, http.StatusBadRequest, "Invalid author id", err)
		return
	}

	var chirps []database.Chirp
	if s != "" {
		chirps, err = cfg.dbQueries.ListChirpsByUserId(r.Context(), id)
	} else {
		chirps, err = cfg.dbQueries.ListChirps(r.Context())
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not list chirps.", err)
		return
	}

	sortQuery := r.URL.Query().Get("sort")
	if sortQuery != "" && sortQuery == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
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

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid chirp UUID.", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not get chirp.", err)
		return
	}

	respondWithJson(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}
