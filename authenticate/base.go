package authenticate

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ZettaAI/auth-server/authorize"
	"github.com/ZettaAI/auth-server/providers"
	"github.com/ZettaAI/auth-server/redis"
	"github.com/ZettaAI/auth-server/utils"
	"github.com/labstack/echo"
)

const (
	// AuthEP main auth endpoint
	AuthEP = "/auth"
	// AuthLogoutEP main logout endpoint
	AuthLogoutEP = "/auth/logout"
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
		return forceLogin(c)
	}
	return validateToken(c, authURL, token)
}

func forceLogin(c echo.Context) error {
	// for convenience redirect users to google login
	// when middle_auth_token is missing
	// happens when a user visits an endpoint directly
	redirectURL := fmt.Sprintf(
		"%v%v%v",
		utils.GetRequestSchemeAndHostURL(c),
		c.Request().Header.Get("X-Forwarded-Prefix"),
		c.Request().Header.Get("X-Forwarded-Uri"),
	)
	c.QueryParams().Set("redirect", redirectURL)
	return providers.GoogleLogin(c)
}

// extractAuthToken helper function
// checks if token is present in
// query params, auth header, or cookie
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
	if val, ok := queryMap[providers.AuthTokenIdentifier]; ok {
		return val[0]
	}

	// check if authorization header has token
	authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	if authHeader != "" {
		authToken := strings.Fields(authHeader)
		if len(authToken) == 2 {
			return authToken[1]
		}
	}

	// check if cookie has token
	token, err := c.Cookie(providers.AuthTokenIdentifier)
	if err == nil {
		return token.Value
	}
	// use none to explicitly indicate missing query param or cookie
	// need to differentiate param set to "" vs missing param
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
			http.StatusBadRequest, fmt.Sprintf("Login at %v to get a token.", authURL))
	}
	// token available, check if present in redis cache
	email, err := redis.GetToken(token)
	if err != nil {
		// bad token, unauthorized 401
		c.Response().Header().Set("WWW-Authenticate", authHeader)
		// discard cookie if exists, for user convenience
		c.SetCookie(&http.Cookie{
			Name:   providers.AuthTokenIdentifier,
			Value:  "",
			MaxAge: -1,
			Path:   "/",
		})
		return c.String(
			http.StatusUnauthorized,
			fmt.Sprintf("Invalid/expired token. Try again."),
		)
	}
	return authorize.Authorize(c, email)
}

// Logout main logout handler.
// User is prompted to login for indentification when visiting /auth/logout
// Email is captured in X-Forwarded-User after successful authentication
// After getting user email, delete associated tokens in redis.
func Logout(c echo.Context) error {
	email := c.Request().Header.Get("X-Forwarded-User")
	res := providers.DeleteUserTokens(email)
	log.Printf("Logged out user: %v, deleted %v keys.", email, res)

	// discard cookie if exists
	c.SetCookie(&http.Cookie{
		Name:   providers.AuthTokenIdentifier,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
	return c.String(http.StatusOK, fmt.Sprintf("%v logged out.", email))
}
