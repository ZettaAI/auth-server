package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
)

var ctx = context.Background()

// Login main login handler.
// AUTHENTICATION:
//   Checks if a given token already exists in redis cache:
//     If yes, user is authenticated, proceed to check authorization.
//     If not, call oauth handler based on x-oauth header (default: Google OAuth)
//       Create a secret token and map it to user email in redis.
// AUTHORIZATION:
func Login(c echo.Context) error {
	// check if token exists in forwarded URL or in headers
	log.Println(c.Request().Header)
	log.Println()

	uri := c.Request().Header.Get("X-Forwarded-Uri")
	queryMap, _ := url.ParseQuery(uri)

	authURL := fmt.Sprintf(
		"%v://%v/auth/google/login", c.Scheme(), c.Request().Host)

	return validateToken(c, authURL, queryMap.Get("middle_auth_token"))
}

func validateToken(c echo.Context, authURL string, token string) error {
	if token == "" {
		// no token, bad request 400
		c.Response().Header().Set(
			"WWW-Authenticate",
			fmt.Sprintf("Bearer realm='%v', error='%v'", authURL, "invalid_token"),
		)
		return c.String(http.StatusBadRequest, "Login required.")
	}

	// token available, check if present in redis cache
	email, err := RedisDB.Get(ctx, token).Result()
	if err != nil {
		// bad token, unauthorized 401
		c.Response().Header().Set(
			"WWW-Authenticate",
			fmt.Sprintf("Bearer realm='%v', error='%v'", authURL, "invalid_token"),
		)
		return c.String(http.StatusUnauthorized, "Invalid/expired token.")
	}
	log.Println(email)
	c.Response().Header().Set("X-Forwarded-User", email)
	return c.String(http.StatusOK, "")
}

// SetToken haha
func SetToken(k string, v string) bool {
	err := RedisDB.Set(ctx, k, v, 0).Err()
	if err != nil {
		panic(err)
	}
	return true
}
