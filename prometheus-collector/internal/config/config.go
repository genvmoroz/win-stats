package config

import (
	"fmt"

	"github.com/genvmoroz/win-stats-prometheus-collector/internal/core"
	"github.com/genvmoroz/win-stats-prometheus-collector/internal/infrastructure"
	"github.com/genvmoroz/win-stats-prometheus-collector/internal/repository/picker"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel    string `envconfig:"APP_LOG_LEVEL" default:"debug"`
	CoreService core.Config
	Infra       infrastructure.Config
	Picker      picker.Config
}

func ReadFromEnv() (Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("process env: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
