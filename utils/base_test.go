package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestGetRequestSchemeAndHostURL(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	req.Header.Set("X-Forwarded-Host", "zetta.ai")
	req.Header.Set(echo.HeaderXForwardedProto, "https")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.Equal(t, "https://zetta.ai", GetRequestSchemeAndHostURL(c))
}
