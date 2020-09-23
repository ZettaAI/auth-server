package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/akhileshh/auth-server/providers"
	"github.com/akhileshh/auth-server/redis"
	"github.com/akhileshh/auth-server/utils"
	"github.com/labstack/echo"
)

const (
	// AuthEP main auth endpoint
	AuthEP = "/auth"
)

// Login main login handler.
// AUTHENTICATION:
//   X-Forwarded-Uri:[/?middle_auth_token=a]
//   Checks if a given token already exists in redis cache:
//     If yes, user is authenticated, proceed to check authorization.
//     If not, call oauth handler based on x-oauth header (default: Google OAuth)
//       Create a secret token and map it to user email in redis.
// AUTHORIZATION: TODO
//   X-Forwarded-Prefix:[/lab2]
func Login(c echo.Context) error {
	// check if token exists in forwarded URL or in headers (TODO)
	log.Printf("Header %v\n", c.Request().Header)
	authURL := utils.GetRequestSchemeAndHostURL(c) + providers.GoogleOAuthLoginEP
	authHeader := fmt.Sprintf("Bearer realm=%v, error=%v", authURL, "invalid_token")
	return validateToken(c, authHeader, extractAuthToken(c))
}

// extractAuthToken helper function to parse query string
func extractAuthToken(c echo.Context) string {
	// check forwarded uri from load balancer/proxy
	// in the form of X-Forwarded-Uri:<string>
	uri := c.Request().Header.Get("X-Forwarded-Uri")
	if uri != "" && strings.IndexRune(uri, '?') != -1 {
		uri = strings.Split(uri, "?")[1]
	} else {
		uri = c.Request().URL.RawQuery
	}
	queryMap, _ := url.ParseQuery(uri)
	return queryMap.Get("middle_auth_token")
}

// validateToken checks for token validity
// if not valid, add header to response to indicate where user can authenticate
// this is currently based on how client handles auth (neuroglancer)
func validateToken(c echo.Context, authHeader string, token string) error {
	// no token provided
	if token == "" {
		// bad request 400
		c.Response().Header().Set("WWW-Authenticate", authHeader)
		return c.String(http.StatusBadRequest, "Login required.")
	}
	// token available, check if present in redis cache
	email, err := redis.GetToken(token)
	if err != nil {
		// bad token, unauthorized 401
		c.Response().Header().Set("WWW-Authenticate", authHeader)
		return c.String(http.StatusUnauthorized, "Invalid/expired token.")
	}
	log.Printf("Logged in user email: %v\n", email)
	// add forward header for backend
	c.Response().Header().Set("X-Forwarded-User", email)
	return c.String(http.StatusOK, "")
}
