package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// New create new echo server with basic configuration.
func New(cfg *Config) *echo.Echo {
	cfg.fillDefaults()

	e := echo.New()

	middlewares := []echo.MiddlewareFunc{}

	middlewares = append(middlewares,
		middleware.CSRFWithConfig(middleware.CSRFConfig{
			TokenLookup:    "cookie:_csrf",
			CookiePath:     "/",
			CookieDomain:   cfg.Domain,
			CookieSecure:   true,
			CookieHTTPOnly: true,
			CookieSameSite: http.SameSiteLaxMode,
		}),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:  cfg.AllowOrigins,
			AllowHeaders:  cfg.AllowHeaders,
			AllowMethods:  cfg.AllowMethods,
			ExposeHeaders: []string{echo.HeaderContentLength},
		}),
		middleware.RecoverWithConfig(middleware.RecoverConfig{
			StackSize: 1 << 10,
			LogLevel:  log.ERROR,
		}),
		middleware.SecureWithConfig(middleware.SecureConfig{
			XSSProtection:         "1; mode=block",
			ContentTypeNosniff:    "nosniff",
			XFrameOptions:         "DENY",
			HSTSMaxAge:            31536000,
			HSTSExcludeSubdomains: true,
			// ContentSecurityPolicy: "default-src 'self'",
		}),
		middleware.Logger(),
	)

	e.Debug = cfg.Debug
	if e.Debug {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.ERROR)
	}

	for _, v := range middlewares {
		e.Use(v)
	}

	e.Server.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Second
	e.Server.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Second
	e.Server.Addr = cfg.Address

	return e
}

// Start starts echo server
func Start(e *echo.Echo, isDevelopment bool) {
	if isDevelopment {
		for _, r := range e.Routes() {
			e.Logger.Infof("%s %s", r.Method, r.Path)
		}

		// Start server
		go func() {
			if err := e.StartServer(e.Server); err != nil {
				if err == http.ErrServerClosed {
					e.Logger.Info("shutting down the server")
				} else {
					e.Logger.Errorf("error shutting down the server: %s", err)
				}
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		e.Logger.Info("received interrupt signal")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	} else {
		// Hide verbose logs and start server normally
		e.HideBanner = true
		e.HidePort = true
		e.Logger.Fatal(e.StartServer(e.Server))
	}
}

// Constants
const (
	gracefulShutdownTimeout = 10 * time.Second
)
