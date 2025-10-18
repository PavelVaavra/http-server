package auth

import (
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
