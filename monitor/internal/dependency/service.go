package dependency

import (
	"github.com/genvmoroz/win-stats-service/internal/core"
	"github.com/genvmoroz/win-stats-service/internal/repository/mem"
	"github.com/genvmoroz/win-stats-service/internal/repository/stats"
	"github.com/genvmoroz/win-stats-service/internal/repository/timegen"
	"github.com/samber/do"
)

func NewService(injector *do.Injector) (*core.Service, error) {
	var (
		timeGenerator = do.MustInvoke[*timegen.TimeGenerator](injector)
		statsRepo     = do.MustInvoke[*stats.Repo](injector)
		memStore      = do.MustInvoke[*mem.Store](injector)
	)

	return core.NewService(timeGenerator, statsRepo, memStore)
}
