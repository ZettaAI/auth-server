package auth

import (
	"fmt"
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
	authURL := utils.GetRequestSchemeAndHostURL(c) + providers.GoogleOAuthLoginEP
	return validateToken(c, authURL, extractAuthToken(c))
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
func validateToken(c echo.Context, authURL string, token string) error {
	authHeader := fmt.Sprintf("Bearer realm=%v, error=%v", authURL, "invalid_token")
	// no token provided
	if token == "" {
		// bad request 400
		c.Response().Header().Set("WWW-Authenticate", authHeader)
		return c.String(
			http.StatusBadRequest, fmt.Sprintf("Please login at %v", authURL))
	}
	// token available, check if present in redis cache
	email, err := redis.GetToken(token)
	if err != nil {
		// bad token, unauthorized 401
		c.Response().Header().Set("WWW-Authenticate", authHeader)
		return c.String(
			http.StatusUnauthorized,
			fmt.Sprintf("Invalid/expired token. Please login at %v", authURL),
		)
	}
	// add forward header for backend
	c.Response().Header().Set("X-Forwarded-User", email)
	return Authorize(c, email)
}
