package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cirius-go/cirius/config"
	"github.com/cirius-go/cirius/pkg/server"
	"github.com/cirius-go/cirius/pkg/util/awsutil"
	"github.com/labstack/echo/v4"
)

var (
	createdAt = time.Now()
	stage     = os.Getenv("STAGE")
	appName   = os.Getenv("APP_NAME")
)

func main() {
	cfg, err := config.Load(&config.LoadOptions{
		Stage:         stage,
		AppName:       appName,
		CustomHandler: loadConfigFromAPSHandler(stage),
	})
	if err != nil {
		panic(err)
	}

	serverCfg := server.Config{
		Debug:        cfg.Debug,
		Address:      fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort),
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
		Domain:       cfg.AppDomain,
		AllowOrigins: cfg.ServerAllowOrigins,
		AllowHeaders: cfg.ServerAllowHeaders,
		AllowMethods: cfg.ServerAllowMethods,
	}

	e := server.New(&serverCfg)

	// init routes
	e.GET("/", checkHealth)

	server.Start(e, serverCfg.Debug)
}

func checkHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"created_at": createdAt,
	})
}

func loadConfigFromAPSHandler(stage string) config.ParameterStoreHandler {
	if stage == "local" {
		return nil
	}

	return func(stage, appName string) error {
		return awsutil.LoadEnvFromAPS(stage, appName, nil)
	}
}
