package providers

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
	"strconv"
	"strings"

	"github.com/ZettaAI/auth-server/utils"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	// GoogleOAuthLoginEP login endpoint
	GoogleOAuthLoginEP = "/auth/google/login"
	// GoogleOAuthCallbackEP google oauth callback endpoint
	GoogleOAuthCallbackEP = "/auth/google/callback"
	googleOAuthUserInfoEP = "https://www.googleapis.com/oauth2/v2/userinfo"
)

var oauthConfig = &oauth2.Config{
	RedirectURL:  fmt.Sprintf("http://localhost:8000%v", GoogleOAuthCallbackEP),
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

// GoogleLogin google oauth login handler
func GoogleLogin(c echo.Context) error {
	// https://developers.google.com/identity/protocols/oauth2/openid-connect#server-flow
	oauthConfig.RedirectURL = utils.GetRequestSchemeAndHostURL(c) + GoogleOAuthCallbackEP
	queryMap := c.QueryParams()
	c.SetCookie(&http.Cookie{
		Name:   "redirectTo",
		Value:  queryMap.Get("redirect"),
		MaxAge: 300,
		Path:   "/",
	})
	return c.Redirect(
		http.StatusTemporaryRedirect,
		oauthConfig.AuthCodeURL(createOauthStateCookie(c)),
	)
}

// GoogleCallback google oauth callback handler
func GoogleCallback(c echo.Context) error {
	oauthConfig.RedirectURL = utils.GetRequestSchemeAndHostURL(c) + GoogleOAuthCallbackEP

	// check if tampered
	oAuthState, err := c.Cookie("oAuthState")
	if err != nil {
		log.Printf("Cookie error: %v", err.Error())
		return c.String(
			http.StatusBadRequest,
			fmt.Sprintf("Could not find cookie: %v", err.Error()),
		)
	}
	if c.FormValue("state") != oAuthState.Value {
		log.Println("Invalid google oauth state.")
		return c.String(http.StatusBadRequest, "Response integrity issue.")
	}

	data, err := getUserInfo(c.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		return c.String(
			http.StatusBadRequest, fmt.Sprintf("Error getting user email: %v", err))
	}

	// got email, generate token and add it to redis cache
	var userInfo map[string]interface{}
	json.Unmarshal(data, &userInfo)
	redirectTo, err := c.Cookie("redirectTo")
	if err != nil {
		log.Printf("Cookie error: %v", err.Error())
		return c.String(
			http.StatusBadRequest,
			fmt.Sprintf("Could not find cookie: %v", err.Error()),
		)
	}

	AddUser(fmt.Sprintf("%v", userInfo["email"]), fmt.Sprintf("%v", userInfo["name"]))
	redirect := redirectTo.Value
	// logged in from neuroglancer
	if strings.Contains(redirect, "appspot.com") {
		token := GetUniqueToken(fmt.Sprintf("%v", userInfo["email"]), false)
		url := fmt.Sprintf("%v?%v=%v", redirect, AuthTokenIdentifier, token)
		return c.Redirect(http.StatusFound, url)
	} else if redirect != "" {
		c.SetCookie(&http.Cookie{
			Name:  AuthTokenIdentifier,
			Value: GetUniqueToken(fmt.Sprintf("%v", userInfo["email"]), true),
			Path:  "/",
			MaxAge: func() int {
				n, err := strconv.Atoi(os.Getenv("AUTH_TOKEN_TEMP_EXPIRY_SECONDS"))
				if err != nil {
					log.Printf("env error AUTH_TOKEN_TEMP_EXPIRY_SECONDS %v", err)
					return 60
				}
				return int(n)
			}(),
		})
		return c.Redirect(http.StatusFound, redirect)
	}
	return c.String(
		http.StatusOK, GetUniqueToken(fmt.Sprintf("%v", userInfo["email"]), false))
}

func createOauthStateCookie(c echo.Context) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie(&http.Cookie{
		Name:   "oAuthState",
		Value:  state,
		MaxAge: 300,
		Path:   "/",
	})
	return state
}

func getUserInfo(code string) ([]byte, error) {
	token, err := oauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		return nil, fmt.Errorf("failed code exchange: %s", err.Error())
	}

	response, err := http.Get(
		fmt.Sprintf("%v?access_token=%v", googleOAuthUserInfoEP, token.AccessToken))
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
