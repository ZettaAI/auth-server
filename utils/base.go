package utils

import (
	"fmt"

	"github.com/labstack/echo"
)

// GetRequestSchemeAndHostURL returns scheme and host as string
// example scheme="https", host="example.com", returns "https://example.com"
func GetRequestSchemeAndHostURL(c echo.Context) string {
	// check headers from load balancer/proxy
	host := c.Request().Header.Get("X-Forwarded-Host")
	scheme := c.Request().Header.Get(echo.HeaderXForwardedProto)

	// if none use host and scheme in request
	if host == "" && scheme == "" {
		host = c.Request().Host
		scheme = c.Scheme()
	}
	return fmt.Sprintf("%v://%v", scheme, host)
}
