package core

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
)

type (
	StatsReporter interface {
		ReportSensorValue(value int64, hardwareName, hardwareType, sensorName, sensorType string)
	}

	StatsProvider interface {
		GetStats(ctx context.Context) (Stats, error)
	}

	Config struct {
		CollectInterval time.Duration `envconfig:"APP_COLLECT_INTERVAL" default:"1s"`
	}

	Service struct {
		cfg           Config
		statsReporter StatsReporter
		statsProvider StatsProvider
	}
)

func NewService(cfg Config, statsReporter StatsReporter, statsProvider StatsProvider) (*Service, error) {
	if lo.IsNil(statsReporter) {
		return nil, fmt.Errorf("stats reporter is nil")
	}
	if lo.IsNil(statsProvider) {
		return nil, fmt.Errorf("stats provider is nil")
	}
	if cfg.CollectInterval <= 0 {
		return nil, fmt.Errorf("collect interval must be greater than 0")
	}
	return &Service{
		cfg:           cfg,
		statsReporter: statsReporter,
		statsProvider: statsProvider,
	}, nil
}

func (s *Service) Collect(ctx context.Context) error {
	ticker := time.NewTicker(s.cfg.CollectInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			stats, err := s.statsProvider.GetStats(ctx)
			if err != nil {
				return fmt.Errorf("get stats: %w", err)
			}

			for _, hardware := range stats.Hardware {
				for _, sensor := range hardware.Sensors {
					s.statsReporter.ReportSensorValue(
						sensor.Value.Value,
						constructHardwareName(hardware),
						hardware.Type,
						constructSensorName(sensor),
						sensor.Type,
					)
				}
			}
		}
	}
}

func constructHardwareName(hardware Hardware) string {
	return fmt.Sprintf("%s (%s)", hardware.Name, hardware.ID)
}

func constructSensorName(sensor Sensor) string {
	return fmt.Sprintf("%s (%s)", sensor.Name, sensor.ID)
}
