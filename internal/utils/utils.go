package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plain password with a bcrypt hash.
// Returns nil if the password matches, otherwise returns an error with a message.
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("invalid password: does not match stored hash")
	}
	return nil
}

// GenerateRandomToken creates a secure random string for CSRF tokens.
func GenerateRandomToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		// fallback to less secure random if needed
		return base64.URLEncoding.EncodeToString([]byte("fallback-csrf-token"))
	}
	return base64.URLEncoding.EncodeToString(b)
}
