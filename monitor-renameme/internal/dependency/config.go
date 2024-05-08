package dependency

import (
	"github.com/genvmoroz/win-stats-service/internal/config"
	"github.com/samber/do"
)

func NewConfig(_ *do.Injector) (config.Config, error) {
	return config.FromEnv()
}
