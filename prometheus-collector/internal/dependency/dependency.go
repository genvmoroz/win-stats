package dependency

import (
	"context"
	"fmt"

	"github.com/genvmoroz/win-stats-prometheus-collector/internal/config"
	"github.com/genvmoroz/win-stats-prometheus-collector/internal/core"
	"github.com/genvmoroz/win-stats-prometheus-collector/internal/infrastructure"
	"github.com/genvmoroz/win-stats-prometheus-collector/internal/repository/picker"
	"github.com/genvmoroz/win-stats-prometheus-collector/internal/repository/prometheus"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type Dependency struct {
	coreService *core.Service
	infraServer *infrastructure.Server
}

func MustBuild(ctx context.Context) *Dependency {
	injector := do.DefaultInjector

	do.Provide(injector, NewConfig)
	do.Provide(injector, NewLogger)
	do.Provide(injector, NewInfraServer)
	do.Provide(injector, NewPrometheusStatsReporter)
	do.Provide(injector, NewStatsPickerRepo(ctx))
	do.Provide(injector, NewCoreService)

	return &Dependency{
		coreService: do.MustInvoke[*core.Service](injector),
		infraServer: do.MustInvoke[*infrastructure.Server](injector),
	}
}

func (d *Dependency) GetCoreService() *core.Service {
	return d.coreService
}

func (d *Dependency) GetInfraServer() *infrastructure.Server {
	return d.infraServer
}

func NewCoreService(injector *do.Injector) (*core.Service, error) {
	var (
		cfg           = do.MustInvoke[config.Config](injector)
		statsReporter = do.MustInvoke[*prometheus.StatsReporter](injector)
		statsProvider = do.MustInvoke[*picker.Repo](injector)
	)
	return core.NewService(cfg.CoreService, statsReporter, statsProvider)
}

func NewConfig(_ *do.Injector) (config.Config, error) {
	cfg, err := config.ReadFromEnv()
	if err != nil {
		return config.Config{}, fmt.Errorf("read config from env: %w", err)
	}
	if err = cfg.Validate(); err != nil {
		return config.Config{}, fmt.Errorf("validate config: %w", err)
	}
	return cfg, nil
}

func NewPrometheusStatsReporter(_ *do.Injector) (*prometheus.StatsReporter, error) {
	reporter := prometheus.NewStatsReporter()
	if err := reporter.Register(); err != nil {
		return nil, fmt.Errorf("register prometheus stats reporter: %w", err)
	}
	return reporter, nil
}

func NewStatsPickerRepo(ctx context.Context) func(injector *do.Injector) (*picker.Repo, error) {
	return func(injector *do.Injector) (*picker.Repo, error) {
		cfg := do.MustInvoke[config.Config](injector)
		return picker.NewRepo(ctx, cfg.Picker)
	}
}

func NewInfraServer(injector *do.Injector) (*infrastructure.Server, error) {
	var (
		cfg    = do.MustInvoke[config.Config](injector)
		logger = do.MustInvoke[logrus.FieldLogger](injector)
	)
	return infrastructure.NewServer(cfg.Infra, logger)
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
