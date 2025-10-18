package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWTValidateJWT(t *testing.T) {
	cases := []struct {
		tokenSecretMake     string
		tokenSecretValidate string
		expiresIn           time.Duration
		errorExpected       bool
		idEqualness         bool
	}{
		{
			tokenSecretMake:     "AllYourBase",
			tokenSecretValidate: "AllYourBase",
			expiresIn:           time.Hour,
			errorExpected:       false,
			idEqualness:         true,
		},
		{
			tokenSecretMake:     "AllYourBase",
			tokenSecretValidate: "allYourBase",
			expiresIn:           time.Hour,
			errorExpected:       true,
			idEqualness:         false,
		},
		{
			tokenSecretMake:     "AllYourBase",
			tokenSecretValidate: "AllYourBase",
			expiresIn:           -time.Hour,
			errorExpected:       true,
			idEqualness:         false,
		},
	}

	for _, c := range cases {
		id := uuid.New()
		tokenString, err := MakeJWT(id, c.tokenSecretMake, c.expiresIn)
		if err != nil {
			t.Errorf("MakeJWT(%v, %v, %v) returns an error %v", id, c.tokenSecretMake, c.expiresIn, err.Error())
		}
		validatedId, err := ValidateJWT(tokenString, c.tokenSecretValidate)
		errActual := (err != nil)
		if errActual != c.errorExpected {
			t.Errorf("%v != %v", errActual, c.errorExpected)
			if err != nil {
				t.Errorf("ValidateJWT(%v, %v) returns an error %v", tokenString, c.tokenSecretValidate, err.Error())
			}
		}
		idActual := (id == validatedId)
		if idActual != c.idEqualness {
			t.Errorf("%v != %v", idActual, c.idEqualness)
		}
	}
}

func TestBearerHeader(t *testing.T) {
	cases := []struct {
		headers           http.Header
		bearerTokenOutput string
		errorExpected     bool
	}{
		{
			headers: http.Header{
				"Authorization": {"  Bearer bearerToken"},
			},
			bearerTokenOutput: "bearerToken",
			errorExpected:     false,
		},
		{
			headers: http.Header{
				"Authorization": {""},
			},
			bearerTokenOutput: "",
			errorExpected:     true,
		},
		{
			headers: http.Header{
				"Accept": {"  Bearer bearerToken"},
			},
			bearerTokenOutput: "",
			errorExpected:     true,
		},
	}

	for _, c := range cases {
		actualToken, err := GetBearerToken(c.headers)
		errActual := (err != nil)
		if errActual != c.errorExpected {
			t.Errorf("%v != %v", errActual, c.errorExpected)
			if err != nil {
				t.Errorf("GetBearerToken(%v) returns an error %v", c.headers, err.Error())
			}
		}
		if actualToken != c.bearerTokenOutput {
			t.Errorf("%v != %v", actualToken, c.bearerTokenOutput)
		}
	}
}
