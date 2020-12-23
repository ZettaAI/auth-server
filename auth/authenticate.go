package authenticate

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/akhileshh/auth-server/authorize"
	"github.com/akhileshh/auth-server/providers"
	"github.com/akhileshh/auth-server/redis"
	"github.com/akhileshh/auth-server/utils"
	"github.com/labstack/echo"
)

const (
	// AuthEP main auth endpoint
	AuthEP = "/auth"
	// AuthLogoutEP main logout endpoint
	AuthLogoutEP = "/auth/logout_all"
)

// Login main login handler.
// X-Forwarded-Uri:[/?middle_auth_token=a]
// Checks if a given token already exists in redis cache:
//   If yes, user is authenticated, proceed to check authorization.
//   If not, call oauth handler based on x-oauth header (default: Google OAuth)
// 	   Create a token and map it to user email in redis.
func Login(c echo.Context) error {
	// TODO check if token exists in forwarded URL or in headers
	authURL := utils.GetRequestSchemeAndHostURL(c) + providers.GoogleOAuthLoginEP
	token := extractAuthToken(c)
	if token == "none" {
		// for convenience redirect users to google login
		// when query param middle_auth_token is missing
		// this is useful when a user visits an endpoint directly
		c.Request().URL.Query().Set("redirect", "none")
		return providers.GoogleLogin(c)
	}
	return validateToken(c, authURL, token)
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
	if val, ok := queryMap["middle_auth_token"]; ok {
		return val[0]
	}
	return "none"
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
	return authorize.Authorize(c, email)
}

// Logout main logout handler.
// User is prompted to login for indentification when visiting /auth/logout_all
// After getting user email, delete associated tokens in redis.
func Logout(c echo.Context) error {
	// TODO check if token exists in forwarded URL or in headers
	authURL := utils.GetRequestSchemeAndHostURL(c) + providers.GoogleOAuthLoginEP
	token := extractAuthToken(c)
	if token == "none" {
		// for convenience redirect users to google login
		// when query param middle_auth_token is missing
		// this is useful when a user visits an endpoint directly
		c.Request().URL.Query().Set("redirect", "none")
		return providers.GoogleLogin(c)
	}
	return validateToken(c, authURL, token)
}
