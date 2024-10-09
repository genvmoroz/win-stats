package config

import (
	"fmt"

	"github.com/genvmoroz/win-stats-picker/internal/http"
	"github.com/genvmoroz/win-stats-picker/internal/repository/stats"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel string `envconfig:"APP_LOG_LEVEL" default:"info"`

	HTTPServer http.Config
	CachedRepo stats.CachedRepoConfig
}

func FromEnv() (Config, error) {
	config := Config{}

	err := envconfig.Process("", &config)
	if err != nil {
		return config, fmt.Errorf("load config: %w", err)
	}

	if err = config.validate(); err != nil {
		return config, fmt.Errorf("validate config: %w", err)
	}

	return config, nil
}

func (c *Config) validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
