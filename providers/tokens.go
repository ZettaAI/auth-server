package providers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ZettaAI/auth-server/redis"
)

const (
	// AuthTokenIdentifier login endpoint
	AuthTokenIdentifier = "middle_auth_token"
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
func GetUniqueToken(email string, temporary bool) string {
	nDays := func() int {
		n, err := strconv.Atoi(os.Getenv("AUTH_TOKEN_EXPIRY_DAYS"))
		if err != nil {
			return 30 // use 30 as default number of days
		}
		return int(n)
	}()

	expiration := time.Hour * time.Duration(nDays*24)
	if temporary {
		// expire quickly when user visits endpoints direclty (not via neuroglancer)
		nSeconds := func() int {
			n, err := strconv.Atoi(os.Getenv("AUTH_TOKEN_TEMP_EXPIRY_SECONDS"))
			if err != nil {
				return 60
			}
			return int(n)
		}()
		expiration = time.Second * time.Duration(nSeconds)
	}

	token, _ := generateRandomString(32)
	for !redis.SetTokenIfNotExists(token, email, expiration) {
		token, _ = generateRandomString(32)
	}

	// for efficient logout keep a copy of token starting with associated email
	// this is more convienient than maintaining a set of tokens with email as key
	// because expiration can't be added to set members
	// when a user wants to logout, simply scan all keys starting with email
	// and delete the combined key and token key
	redis.SetToken(
		fmt.Sprintf("%v:%v", email, token), "", expiration)
	return token
}

// DeleteUserTokens delete all tokens of given user
// TODO use goroutine for faster response
func DeleteUserTokens(email string) int64 {
	// gather all keys that need to be deleted
	userCombinedTokens := redis.GetTokensStartingWith(email)
	userTokens := make([]string, len(userCombinedTokens))
	for i, key := range userCombinedTokens {
		// combined token has the format email:token
		userTokens[i] = strings.Split(key, ":")[1]
	}
	// delete all with one redis call for efficiency
	return redis.DeleteTokens(append(userCombinedTokens, userTokens...)...)
}
