package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cirius-go/cirius/pkg/server"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewEchoServer(t *testing.T) {
	// Create a new configuration with some custom values for testing
	cfg := &server.Config{
		Address:      ":8080",
		Debug:        true,
		ReadTimeout:  10,
		WriteTimeout: 20,
		AllowOrigins: []string{"*"},
	}

	// Create a new echo server with the custom configuration
	e := server.New(cfg)

	// Test that the server is configured correctly
	assert.Equal(t, e.Debug, true)
	assert.Equal(t, e.Server.ReadTimeout, 10*time.Second)
	assert.Equal(t, e.Server.WriteTimeout, 20*time.Second)
	assert.Equal(t, e.Server.Addr, ":8080")

	// Create a new request and recorder for testing middleware
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Test that the CSRF middleware is configured correctly
	e.ServeHTTP(rec, req)
	assert.Equal(t, rec.Header().Get(echo.HeaderSetCookie), "_csrf=")
	assert.Equal(t, rec.Result().StatusCode, http.StatusForbidden)

	// Test that the CORS middleware is configured correctly
	req.Header.Set(echo.HeaderOrigin, "example.com")
	e.ServeHTTP(rec, req)
	assert.Equal(t, rec.Header().Get(echo.HeaderAccessControlAllowOrigin), "example.com")
	assert.Equal(t, rec.Result().StatusCode, http.StatusForbidden)

	// Test that the recover middleware is configured correctly
	// req = httptest.NewRequest(http.MethodGet, "/", nil)
	// rec = httptest.NewRecorder()
	// h := func(c echo.Context) error {
	// 	panic("test panic")
	// }
	// e.ServeHTTP(rec, req.WithContext(e.NewContext(req, rec)))
	// assert.Equal(t, rec.Result().StatusCode, http.StatusInternalServerError)

	// Test that the secure middleware is configured correctly
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, rec.Header().Get(echo.HeaderXFrameOptions), "DENY")
	assert.Equal(t, rec.Result().StatusCode, http.StatusOK)
}
