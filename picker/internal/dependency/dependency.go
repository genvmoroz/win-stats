package dependency

import (
	"fmt"

	"github.com/genvmoroz/win-stats-picker/internal/config"
	"github.com/genvmoroz/win-stats-picker/internal/core"
	"github.com/genvmoroz/win-stats-picker/internal/http"
	"github.com/genvmoroz/win-stats-picker/internal/repository/stats"
	"github.com/genvmoroz/win-stats-picker/internal/repository/timegen"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type Dependency struct {
	httpServer *http.Server
}

func MustBuild() Dependency {
	injector := do.DefaultInjector

	do.ProvideValue(injector, timegen.NewTimeGenerator())

	do.Provide(injector, NewConfig)
	do.Provide(injector, NewLogger)
	do.Provide(injector, NewStatsRepo)
	do.Provide(injector, NewCoreService)
	do.Provide(injector, NewRouter)
	do.Provide(injector, NewHTTPServer)

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

func NewHTTPServer(injector *do.Injector) (*http.Server, error) {
	var (
		cfg    = do.MustInvoke[config.Config](injector)
		router = do.MustInvoke[*http.Router](injector)
		logger = do.MustInvoke[logrus.FieldLogger](injector)
	)

	return http.NewServer(cfg.HTTPServer, router, logger)
}

func NewCoreService(injector *do.Injector) (*core.Service, error) {
	statsRepo := do.MustInvoke[*stats.Repo](injector)

	return core.NewService(statsRepo)
}
