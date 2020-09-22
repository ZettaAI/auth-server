package main

import (
	"net/http"

	"github.com/akhileshh/auth-server/auth"
	"github.com/akhileshh/auth-server/handlers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET(auth.AuthEP, auth.Login)
	e.GET(handlers.GoogleOAuthLoginEP, handlers.GoogleLogin)
	e.GET(handlers.GoogleOAuthCallbackEP, handlers.GoogleCallback)

	e.Logger.Fatal(e.Start(":8000"))
}
