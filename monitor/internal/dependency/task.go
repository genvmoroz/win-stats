package dependency

import (
	"github.com/genvmoroz/win-stats-service/internal/config"
	"github.com/genvmoroz/win-stats-service/internal/core/autocleanup"
	"github.com/genvmoroz/win-stats-service/internal/repository/mem"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

func NewAutoCleanup(injector *do.Injector) (*autocleanup.Task, error) {
	var (
		cfg    = do.MustInvoke[config.Config](injector)
		store  = do.MustInvoke[*mem.Store](injector)
		logger = do.MustInvoke[logrus.FieldLogger](injector)
	)

	return autocleanup.NewTask(cfg.AutoCleanupTask, store, logger)
}
