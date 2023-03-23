package config

import (
	"errors"
	"fmt"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

// ParameterStoreHandler is used to get parameter from another service and must
// be expose all parameters to the env.
type ParameterStoreHandler func(stage, appName string) error

// Config represents the config of application.
type Config struct {
	Debug     bool   `env:"DEBUG"`
	AppDomain string `env:"APP_DOMAIN"`

	ServerHost         string `env:"SERVER_HOST"`
	ServerPort         int    `env:"SERVER_PORT"`
	ServerReadTimeout  int    `env:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout int    `env:"SERVER_WRITE_TIMEOUT"`

	ServerAllowOrigins []string `env:"SERVER_ALLOW_ORIGINS"`
	ServerAllowHeaders []string `env:"SERVER_ALLOW_HEADERS"`
	ServerAllowMethods []string `env:"SERVER_ALLOW_METHODS"`
}

// Validate the config.
func (c *Config) Validate() error {
	return nil
}

// LoadOptions is options for function load.
type LoadOptions struct {
	Stage         string
	AppName       string
	CustomHandler ParameterStoreHandler
}

func (o *LoadOptions) Validate() error {
	if o.AppName == "" {
		return errors.New("`App Name` must be defined for configuration.")
	}

	return nil
}

// Load will loads all variables from .env file & .env.{stage} file and the
// custom parameter store's variables to env process. After that, It will parse
// to a Config.
//
// Stage must be expose to env process if application not runs in the local.
func Load(opts *LoadOptions) (*Config, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	stage := opts.Stage

	if stage == "" {
		stage = "local"
	}

	supportEnvFiles := []string{".env", fmt.Sprintf(".env.%s", stage)}

	if err := godotenv.Load(supportEnvFiles...); err != nil {
		return nil, err
	}

	if opts.CustomHandler != nil {
		if err := opts.CustomHandler(stage, opts.AppName); err != nil {
			return nil, err
		}
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
