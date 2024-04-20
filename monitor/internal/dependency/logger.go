package dependency

import (
	"fmt"

	"github.com/genvmoroz/win-stats-service/internal/config"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

func NewLogger(injector *do.Injector) (logrus.FieldLogger, error) {
	cfg := do.MustInvoke[config.Config](injector)

	logger := logrus.New()
	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("parser log level: %w", err)
	}
	logger.SetLevel(lvl)

	return logger, nil
}
