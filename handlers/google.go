package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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

var oauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

// GoogleLogin google oauth login handler
func GoogleLogin(c echo.Context) error {
	// https://developers.google.com/identity/protocols/oauth2/openid-connect#server-flow
	oauthConfig.RedirectURL = fmt.Sprintf(
		"%v://%v/auth/google/callback", c.Scheme(), c.Request().Host)
	return c.Redirect(
		http.StatusTemporaryRedirect,
		oauthConfig.AuthCodeURL(createOauthStateCookie(c)),
	)
}

func createOauthStateCookie(c echo.Context) string {
	expiration := time.Now().Add(7 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oAuthState", Value: state, Expires: expiration}
	c.SetCookie(&cookie)
	return state
}

// GoogleOAuthCallback google oauth callback handler
func GoogleOAuthCallback(c echo.Context) error {
	oAuthState, _ := c.Cookie("oAuthState")
	if c.FormValue("state") != oAuthState.Value {
		log.Println("Invalid google oauth state.")
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	data, err := getUserInfo(c.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	fmt.Fprintf(c.Response().Writer, "UserInfo: %s\n", data)
	return nil
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
