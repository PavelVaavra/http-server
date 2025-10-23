package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Use jwt.NewWithClaims
// Use jwt.SigningMethodHS256 as the signing method.
// Use jwt.RegisteredClaims as the claims.
// Set the Issuer to "chirpy"
// Set IssuedAt to the current time in UTC
// Set ExpiresAt to the current time plus the expiration time (expiresIn)
// Set the Subject to a stringified version of the user's id
// Use token.SignedString to sign the token with the secret key.
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

// Use the jwt.ParseWithClaims function to validate the signature of the JWT and extract the claims into a *jwt.Token struct.
// An error will be returned if the token is invalid or has expired.
// If all is well with the token, use the token.Claims interface to get access to the user's id from the claims
// (which should be stored in the Subject field). Return the id as a uuid.UUID.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

// Auth information will come into our server in the Authorization header. Its value will look like this: Bearer TOKEN_STRING
// This function should look for the Authorization header in the headers parameter and return the TOKEN_STRING if it exists
// (stripping off the Bearer prefix and whitespace). If the header doesn't exist, return an error. This is an easy one to write a unit test for,
// and I'd recommend doing so.
func GetBearerToken(headers http.Header) (string, error) {
	rawValue := headers.Get("Authorization")
	if rawValue == "" {
		return "", fmt.Errorf("no authorization header")
	}
	rawValue = strings.TrimSpace(rawValue)
	parts := strings.SplitN(rawValue, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("header value not in format \"Bearer TOKEN_STRING\"")
	}
	return parts[1], nil
}

// Authorization: ApiKey THE_KEY_HERE
func GetAPIKey(headers http.Header) (string, error) {
	rawValue := headers.Get("Authorization")
	if rawValue == "" {
		return "", fmt.Errorf("no authorization header")
	}
	rawValue = strings.TrimSpace(rawValue)
	parts := strings.SplitN(rawValue, " ", 2)
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", fmt.Errorf("header value not in format \"ApiKey THE_KEY_HERE\"")
	}
	return parts[1], nil
}
