package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"
)

type StatsRepo interface {
	GetSensorsByHardware(ctx context.Context) (map[Hardware][]Sensor, error)
}

type Service struct {
	statsRepo StatsRepo
}

func NewService(statsRepo StatsRepo) (*Service, error) {
	if lo.IsNil(statsRepo) {
		return nil, errors.New("stats repo is nil")
	}
	return &Service{
		statsRepo: statsRepo,
	}, nil
}

func (s *Service) GetStats(ctx context.Context) (GetStatsResponse, error) {
	sensorsByHardware, err := s.statsRepo.GetSensorsByHardware(ctx)
	if err != nil {
		return GetStatsResponse{}, fmt.Errorf("get sensors: %w", err)
	}

	return GetStatsResponse{
		Stats: sensorsByHardware,
	}, nil
}
