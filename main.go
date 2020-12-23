package main

import (
	"net/http"

	"github.com/akhileshh/auth-server/authenticate"
	"github.com/akhileshh/auth-server/providers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var corsConfig = middleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodPost,
		http.MethodDelete,
	},
	AllowHeaders: []string{echo.HeaderAuthorization},
	ExposeHeaders: []string{
		echo.HeaderContentType,
		echo.HeaderContentLength,
		echo.HeaderContentEncoding,
		echo.HeaderAccept,
		echo.HeaderWWWAuthenticate,
	},
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(corsConfig))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET(authenticate.AuthEP, authenticate.Login)
	e.GET(authenticate.AuthLogoutEP, authenticate.Logout)
	e.GET(providers.GoogleOAuthLoginEP, providers.GoogleLogin)
	e.GET(providers.GoogleOAuthCallbackEP, providers.GoogleCallback)

	e.Logger.Fatal(e.Start(":8000"))
}
