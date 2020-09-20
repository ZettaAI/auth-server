package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akhileshh/auth-server/handlers"
	"github.com/labstack/echo"
)

/*
1. user sends token
2. check if token exists in redis (same cluster but different release)
	* if yes check permissions
	* if no, get user info, create a new token and map to email in redis
3. permissions
	* https://echo.labstack.com/middleware/casbin-auth
*/

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/auth/google/login", handlers.Login)
	e.GET("/auth/google/callback", handlers.OAuthCallback)

	SetToken("a", "akhileshhalageri@gmail.com")
	log.Println(fmt.Sprintf("a: %v", GetToken("a")))

	e.Logger.Fatal(e.Start(":8000"))
}
