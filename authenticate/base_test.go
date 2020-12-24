package authenticate

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ZettaAI/auth-server/providers"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestExtractAuthToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	res := extractAuthToken(e.NewContext(req, rec))
	assert.Equal(t, "none", res)

	token := "abc123"
	req.Header.Set(
		"X-Forwarded-Uri",
		fmt.Sprintf("path?%v=%v", providers.AuthTokenIdentifier, token),
	)
	res = extractAuthToken(e.NewContext(req, rec))
	assert.Equal(t, token, res)
}

// TODO Fix test (needs redis)
// func TestLogin(t *testing.T) {
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, AuthEP, strings.NewReader(""))
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	Login(c)
// 	assert.Equal(t, http.StatusBadRequest, c.Response().Status)

// 	token := "abc123"
// 	req.Header.Set(
// 		"X-Forwarded-Uri",
// 		fmt.Sprintf("path?%v=%v", providers.AuthTokenIdentifier, token),
// 	)
// 	c = e.NewContext(req, rec)
// 	Login(c)
// 	assert.Equal(t, http.StatusUnauthorized, c.Response().Status)
// }
