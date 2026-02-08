// Package core contains domain logic for collecting and reporting stats.
package core

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type (
	StatsReporter interface {
		ReportSensorValue(value int64, host, hardwareID, hardwareName, hardwareType, sensorID, sensorName, sensorType string)
	}

	StatsProvider interface {
		GetStats(ctx context.Context) (Stats, error)
	}

	Config struct {
		CollectInterval time.Duration `envconfig:"APP_COLLECT_INTERVAL" default:"1s"`
		CollectTimeout  time.Duration `envconfig:"APP_COLLECT_TIMEOUT" default:"10s"`
		CollectAttempts uint          `envconfig:"APP_COLLECT_ATTEMPTS" default:"2"`
	}

	Service struct {
		cfg               Config
		logger            *zap.SugaredLogger
		statsReporter     StatsReporter
		statsProviders    map[string]StatsProvider
		singleflightGroup singleflight.Group
	}
)

func NewService(cfg Config, logger *zap.SugaredLogger, statsReporter StatsReporter, statsProviders map[string]StatsProvider) (*Service, error) {
	if lo.IsNil(statsReporter) {
		return nil, fmt.Errorf("stats reporter is nil")
	}
	if lo.IsNil(logger) {
		return nil, fmt.Errorf("logger is nil")
	}
	if len(statsProviders) == 0 {
		return nil, fmt.Errorf("stats providers list is empty")
	}
	if cfg.CollectInterval <= 0 {
		return nil, fmt.Errorf("collect interval must be greater than 0")
	}
	if cfg.CollectTimeout <= 0 {
		return nil, fmt.Errorf("collect timeout must be greater than 0")
	}
	if cfg.CollectAttempts == 0 {
		return nil, fmt.Errorf("collect attempts must be greater than 0")
	}
	return &Service{
		cfg:               cfg,
		logger:            logger,
		statsReporter:     statsReporter,
		statsProviders:    statsProviders,
		singleflightGroup: singleflight.Group{},
	}, nil
}

func (s *Service) Collect(ctx context.Context) error {
	ticker := time.NewTicker(s.cfg.CollectInterval)
	defer ticker.Stop()

	// collect stats immediately to avoid waiting for the first tick
	s.collectStatsFromAllProvidersWithRetries(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			s.collectStatsFromAllProvidersWithRetries(ctx)
		}
	}
}

func (s *Service) collectStatsFromAllProvidersWithRetries(ctx context.Context) {
	for host, statsProvider := range s.statsProviders {
		select {
		case <-ctx.Done():
			return
		default:
		}

		go func() {
			_, err, _ := s.singleflightGroup.Do(
				host,
				func() (any, error) {
					s.collectStatsFromOneProviderWithRetries(ctx, host, statsProvider)
					return nil, nil
				},
			)
			if err != nil {
				s.logger.Errorf("unexpected error returned from singleflight: %s", err.Error())
			}
		}()
	}
}

func (s *Service) collectStatsFromOneProviderWithRetries(ctx context.Context, host string, statsProvider StatsProvider) {
	start := time.Now()

	reqCtx, cancel := context.WithTimeout(ctx, s.cfg.CollectTimeout)
	defer cancel()

	err := retry.Do(
		func() error {
			return s.collectStatsFromOneProvider(reqCtx, host, statsProvider)
		},
		retry.Attempts(s.cfg.CollectAttempts),
		retry.Context(reqCtx),
	)
	if err != nil {
		s.logger.With(
			"host", host,
			"duration", time.Since(start),
		).Errorf("successive retries to collect stats failed: %v", err)
	}
}

func (s *Service) collectStatsFromOneProvider(ctx context.Context, host string, statsProvider StatsProvider) error {
	stats, err := statsProvider.GetStats(ctx)
	if err != nil {
		return fmt.Errorf("get stats: %w", err)
	}

	for _, hardware := range stats.Hardware {
		for _, sensor := range hardware.Sensors {
			s.statsReporter.ReportSensorValue(
				sensor.Value.Value,
				host,
				hardware.ID,
				hardware.Name,
				hardware.Type,
				sensor.ID,
				sensor.Name,
				sensor.Type,
			)
		}
	}

	return nil
}
