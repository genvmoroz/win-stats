package dependency

import (
	"github.com/genvmoroz/custom-collector/internal/core"
	"github.com/genvmoroz/custom-collector/internal/repository/mem"
	"github.com/genvmoroz/custom-collector/internal/repository/stats"
	"github.com/genvmoroz/custom-collector/internal/repository/timegen"
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
