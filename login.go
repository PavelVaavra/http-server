package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/PavelVaavra/http-server/internal/auth"
	"github.com/PavelVaavra/http-server/internal/database"
)

// Update the POST /api/login endpoint to return a refresh token, as well as an access token:
// Access tokens (JWTs) should expire after 1 hour. Expiration time is stored in the exp claim. You can remove the optional expires_in_seconds parameter
// from the endpoint.
// Refresh tokens should expire after 60 days. Expiration time is stored in the database.
// The revoked_at field should be null when the token is created.
// {
//   "id": "5a47789c-a617-444a-8a80-b50359247804",
//   "created_at": "2021-07-01T00:00:00Z",
//   "updated_at": "2021-07-01T00:00:00Z",
//   "email": "lane@example.com",
//   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
//   "refresh_token": "56aa826d22baab4b5ec2cea41a59ecbba03e542aedbb31d9b80326ac8ffcfa2a"
// }

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type userParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	tokenString, err := auth.MakeJWT(user.ID, cfg.tokenSecret, 3600*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot create JWT", err)
		return
	}

	refreshToken, _ := auth.MakeRefreshToken()

	_, err = cfg.dbQueries.RefreshToken(r.Context(), database.RefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create entry in refresh_token table", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        tokenString,
		RefreshToken: refreshToken,
	})
}
