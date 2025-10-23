package main

import (
	"net/http"

	"github.com/PavelVaavra/http-server/internal/auth"
	"github.com/google/uuid"
)

// 1. Add a new DELETE /api/chirps/{chirpID} route to your server that deletes a chirp from the database by its id.
// - This is an authenticated endpoint, so be sure to check the token in the header. Only allow the deletion of a chirp if the user is the author of the chirp.
// - If they are not, return a 403 status code.
// 2. If the chirp is deleted successfully, return a 204 status code.
// 3. If the chirp is not found, return a 404 status code.
func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
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

	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid chirp UUID.", err)
		return
	}

	chirp, err := cfg.dbQueries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not get chirp", err)
		return
	}

	if chirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "Chirp wasn't created by the user who tries to delete it", err)
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Chirp could not be deleted", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
