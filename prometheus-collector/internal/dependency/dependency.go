// Package dependency wires the application dependencies using a simple DI container.
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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	do.Provide(injector, NewStatsPickerRepos(ctx))
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
		cfg            = do.MustInvoke[config.Config](injector)
		logger         = do.MustInvoke[*zap.SugaredLogger](injector)
		statsReporter  = do.MustInvoke[*prometheus.StatsReporter](injector)
		statsProviders = do.MustInvoke[map[string]core.StatsProvider](injector)
	)
	return core.NewService(cfg.CoreService, logger, statsReporter, statsProviders)
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

func NewStatsPickerRepos(ctx context.Context) func(injector *do.Injector) (map[string]core.StatsProvider, error) {
	return func(injector *do.Injector) (map[string]core.StatsProvider, error) {
		logger := do.MustInvoke[*zap.SugaredLogger](injector)
		cfg := do.MustInvoke[config.Config](injector)

		providers := make(map[string]core.StatsProvider, len(cfg.PickerHosts))
		for _, host := range cfg.PickerHosts {
			repo, err := picker.NewRepo(ctx, host)
			if err != nil {
				return nil, fmt.Errorf("init picker repo for host %s: %w", host, err)
			}
			if err = repo.HealthCheck(ctx); err != nil {
				// it's ok to continue if health check fails,
				// it's not a critical problem if one of the hosts is down,
				// it will be reported in the next collection cycle
				logger.Errorf("health check failed for host %s: %v", host, err)
			}
			providers[host] = repo
		}

		return providers, nil
	}
}

func NewInfraServer(injector *do.Injector) (*infrastructure.Server, error) {
	var (
		cfg    = do.MustInvoke[config.Config](injector)
		logger = do.MustInvoke[*zap.SugaredLogger](injector)
	)
	return infrastructure.NewServer(cfg.Infra, logger)
}

func NewLogger(injector *do.Injector) (*zap.SugaredLogger, error) {
	cfg := do.MustInvoke[config.Config](injector)

	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	base, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("build zap logger: %w", err)
	}
	return base.Sugar(), nil
}
