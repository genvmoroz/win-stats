package config

import (
	"fmt"
	"time"

	"github.com/genvmoroz/win-stats-prometheus-collector/internal/core"
	"github.com/genvmoroz/win-stats-prometheus-collector/internal/infrastructure"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel       string        `envconfig:"APP_LOG_LEVEL" default:"debug"`
	CollectTimeout time.Duration `envconfig:"APP_COLLECT_TIMEOUT" default:"10s"`
	PickerHosts    []string      `envconfig:"APP_PICKER_HOSTS" validate:"required"`

	CoreService core.Config
	Infra       infrastructure.Config
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
