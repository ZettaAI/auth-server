package main

import (
	"log"
	"net/http"

	"github.com/akhileshh/auth-server/handlers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		log.Println(c.Request(), c.RealIP())
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/auth/google/login", handlers.Login)
	e.GET("/auth/google/callback", handlers.OAuthCallback)

	e.Logger.Fatal(e.Start(":8000"))
}
