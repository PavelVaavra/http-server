package main

import (
	"net/http"
	"time"

	"github.com/PavelVaavra/http-server/internal/auth"
)

// Create a POST /api/refresh endpoint. This new endpoint does not accept a request body, but does require a refresh token to be present in the headers,
//  in the same Authorization: Bearer <token> format.
// Look up the token in the database. If it doesn't exist, or if it's expired, respond with a 401 status code. Otherwise, respond with a 200 code and
// this shape:

// {
//   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
// }

// The token field should be a newly created access token for the given user that expires in 1 hour. I wrote a GetUserFromRefreshToken SQL query.

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No refresh token in the header", err)
		return
	}

	refreshTokenEntry, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No refresh token in the db table", err)
		return
	}

	if time.Now().After(refreshTokenEntry.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", err)
		return
	}

	if refreshTokenEntry.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", err)
		return
	}

	tokenString, err := auth.MakeJWT(refreshTokenEntry.UserID, cfg.tokenSecret, 3600*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot create JWT", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		Token: tokenString,
	})
}
