package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestExtractAuthToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	res, _ := extractAuthToken(e.NewContext(req, rec))
	assert.Equal(t, "", res)

	token := "abc123"
	req.Header.Set("X-Forwarded-Uri", fmt.Sprintf("path?middle_auth_token=%v", token))
	res, _ = extractAuthToken(e.NewContext(req, rec))
	assert.Equal(t, token, res)
}

func TestLogin(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, AuthEP, strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	Login(c)
	assert.Equal(t, http.StatusBadRequest, c.Response().Status)

	token := "abc123"
	req.Header.Set("X-Forwarded-Uri", fmt.Sprintf("path?middle_auth_token=%v", token))
	c = e.NewContext(req, rec)
	Login(c)
	assert.Equal(t, http.StatusUnauthorized, c.Response().Status)
}
