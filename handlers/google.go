package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	// GoogleOAuthLoginEP login endpoint
	GoogleOAuthLoginEP = "/auth/google/login"
	// GoogleOAuthCallbackEP google oauth callback endpoint
	GoogleOAuthCallbackEP = "/auth/google/callback"
)

var oauthConfig = &oauth2.Config{
	RedirectURL:  fmt.Sprintf("http://localhost:8000%v", GoogleOAuthCallbackEP),
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

// GoogleLogin google oauth login handler
func GoogleLogin(c echo.Context) error {
	// https://developers.google.com/identity/protocols/oauth2/openid-connect#server-flow

	// check headers from load balancer/proxy
	host := c.Request().Header.Get("X-Forwarded-Host")
	scheme := c.Request().Header.Get("X-Forwarded-Proto")

	// if none use host and scheme in request
	if host == "" && scheme == "" {
		host = c.Request().Host
		scheme = c.Scheme()
	}

	oauthConfig.RedirectURL = fmt.Sprintf(
		"%v://%v%v", scheme, host, GoogleOAuthCallbackEP)
	log.Printf("Google callback URL %v\n", oauthConfig.RedirectURL)
	return c.Redirect(
		http.StatusTemporaryRedirect,
		oauthConfig.AuthCodeURL(createOauthStateCookie(c)),
	)
}

// GoogleCallback google oauth callback handler
func GoogleCallback(c echo.Context) error {
	oAuthState, _ := c.Cookie("oAuthState")
	if c.FormValue("state") != oAuthState.Value {
		log.Println("Invalid google oauth state.")
		return c.String(http.StatusBadRequest, "Response integrity issue.")
	}

	data, err := getUserInfo(c.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		return c.String(http.StatusBadRequest, "Error getting user email.")
	}

	// got email, generate token and add it to redis cache
	var userInfo map[string]interface{}
	json.Unmarshal(data, &userInfo)
	token := GetUniqueToken(fmt.Sprintf("%v", userInfo["email"]))
	redirectTo := fmt.Sprintf(
		"%v://%v/auth?middle_auth_token=%v", c.Scheme(), c.Request().Host, token)

	log.Printf("Got user email: %v, token: %v", userInfo["email"], token)
	log.Println(redirectTo)
	return c.Redirect(
		http.StatusFound,
		redirectTo,
	)
}

func createOauthStateCookie(c echo.Context) string {
	expiration := time.Now().Add(1 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oAuthState", Value: state, Expires: expiration}
	c.SetCookie(&cookie)
	return state
}

func getUserInfo(code string) ([]byte, error) {
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	response, err := http.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
