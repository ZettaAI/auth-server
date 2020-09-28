package providers

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/akhileshh/auth-server/redis"
)

// https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// generateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GetUniqueToken generates a token not already in cache.
// Tries to create a new token until it is unique.
// Adds the token to cache with user email as value.
func GetUniqueToken(email string) string {
	token, _ := generateRandomString(32)
	for !redis.SetTokenIfNotExists(token, email, 7*24*time.Hour) {
		token, _ = generateRandomString(32)
	}
	return token
}
