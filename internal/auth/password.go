package auth

import "github.com/alexedwards/argon2id"

// hash the password using the argon2id.CreateHash function
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, err
}

// use the argon2id.ComparePasswordAndHash function to compare the password that the user entered in the HTTP request with the password
// that is stored in the database
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
