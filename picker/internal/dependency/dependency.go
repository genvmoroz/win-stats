package dependency

import (
	"context"
	"fmt"

	"github.com/genvmoroz/win-stats/picker/internal/config"
	"github.com/genvmoroz/win-stats/picker/internal/core"
	"github.com/genvmoroz/win-stats/picker/internal/http"
	"github.com/genvmoroz/win-stats/picker/internal/repository/stats"
	"github.com/genvmoroz/win-stats/picker/internal/repository/timegen"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type Dependency struct {
	httpServer *http.Server
}

func MustBuild(ctx context.Context) Dependency {
	injector := do.DefaultInjector

	do.ProvideValue(injector, timegen.NewTimeGenerator())

	do.Provide(injector, NewConfig)
	do.Provide(injector, NewLogger)
	do.Provide(injector, NewStatsRepo)
	do.Provide(injector, NewSingleflightStatsRepo)
	do.Provide(injector, NewCachedStatsRepo)
	do.Provide(injector, NewCoreService)
	do.Provide(injector, NewRouter)
	do.Provide(injector, NewHTTPServer(ctx))

	return Dependency{
		httpServer: do.MustInvoke[*http.Server](injector),
	}
}

func (d *Dependency) HTTPServer() *http.Server {
	return d.httpServer
}

func NewStatsRepo(injector *do.Injector) (*stats.Repo, error) {
	timeGenerator := do.MustInvoke[*timegen.TimeGenerator](injector)

	return stats.NewRepo(timeGenerator)
}

func NewSingleflightStatsRepo(injector *do.Injector) (*stats.SingleflightRepo, error) {
	baseRepo := do.MustInvoke[*stats.Repo](injector)

	return stats.NewSingleflightRepo(baseRepo)
}

func NewCachedStatsRepo(injector *do.Injector) (*stats.CachedRepo, error) {
	var (
		cfg              = do.MustInvoke[config.Config](injector)
		singleflightRepo = do.MustInvoke[*stats.SingleflightRepo](injector)
	)

	return stats.NewCachedRepo(singleflightRepo, cfg.CachedRepo)
}

func NewRouter(injector *do.Injector) (*http.Router, error) {
	service := do.MustInvoke[*core.Service](injector)

	return http.NewRouter(service)
}

func NewConfig(_ *do.Injector) (config.Config, error) {
	return config.FromEnv()
}

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

func NewHTTPServer(ctx context.Context) func(injector *do.Injector) (*http.Server, error) {
	return func(injector *do.Injector) (*http.Server, error) {
		var (
			cfg    = do.MustInvoke[config.Config](injector)
			router = do.MustInvoke[*http.Router](injector)
			logger = do.MustInvoke[logrus.FieldLogger](injector)
		)

		return http.NewServer(ctx, cfg.HTTPServer, router, logger)
	}
}

func NewCoreService(injector *do.Injector) (*core.Service, error) {
	cachedStatsRepo := do.MustInvoke[*stats.CachedRepo](injector)

	return core.NewService(cachedStatsRepo)
}
