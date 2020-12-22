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
	"net/url"
	"os"
	"time"

	"github.com/akhileshh/auth-server/utils"
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
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

// GoogleLogin google oauth login handler
func GoogleLogin(c echo.Context) error {
	// https://developers.google.com/identity/protocols/oauth2/openid-connect#server-flow
	oauthConfig.RedirectURL = utils.GetRequestSchemeAndHostURL(c) + GoogleOAuthCallbackEP
	queryMap, _ := url.ParseQuery(c.Request().URL.RawQuery)

	c.SetCookie(&http.Cookie{
		Name:    "redirectTo",
		Value:   queryMap.Get("redirect"),
		Expires: time.Now().Add(1 * time.Hour),
	})
	return c.Redirect(
		http.StatusTemporaryRedirect,
		oauthConfig.AuthCodeURL(createOauthStateCookie(c)),
	)
}

// GoogleCallback google oauth callback handler
func GoogleCallback(c echo.Context) error {
	// verify if tampered
	oauthConfig.RedirectURL = utils.GetRequestSchemeAndHostURL(c) + GoogleOAuthCallbackEP
	oAuthState, _ := c.Cookie("oAuthState")
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

	redirectTo, _ := c.Cookie("redirectTo")
	if redirectTo.Value == "none" {
		return c.String(http.StatusOK, "")
	}

	token := GetUniqueToken(fmt.Sprintf("%v", userInfo["email"]))
	if redirectTo.Value == "" {
		return c.String(http.StatusOK, token)
	}

	redirectURL := fmt.Sprintf("%v?middle_auth_token=%v", redirectTo.Value, token)
	return c.Redirect(
		http.StatusFound,
		redirectURL,
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
	token, err := oauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed code exchange: %s", err.Error())
	}

	response, err := http.Get(
		fmt.Sprintf("%v?access_token=%v", googleOAuthUserInfoEP, token.AccessToken))
	if err != nil {
		return nil, fmt.Errorf("Failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed read response: %s", err.Error())
	}
	return contents, nil
}
