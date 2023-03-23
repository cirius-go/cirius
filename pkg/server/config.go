package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Config is the server config.
type Config struct {
	Address      string
	ReadTimeout  int
	WriteTimeout int
	Domain       string
	Debug        bool
	AllowOrigins []string
	AllowHeaders []string
	AllowMethods []string
}

// TODO: setup fillDefaults function for server config.
func (c *Config) fillDefaults() {
	if c.Domain == "" {
		c.Domain = ""
	}

	if len(c.AllowOrigins) == 0 {
		c.AllowOrigins = []string{}
	}

	if len(c.AllowHeaders) == 0 {
		c.AllowHeaders = []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization}
	}

	if len(c.AllowMethods) == 0 {
		c.AllowMethods = []string{http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		}
	}

	if c.ReadTimeout == 0 {
		c.ReadTimeout = 15
	}

	if c.WriteTimeout == 0 {
		c.WriteTimeout = 15
	}
}
