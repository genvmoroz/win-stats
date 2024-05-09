package core

//go:generate mockgen -destination=mock/deps.go -package=mock -source=core.go TimeGenerator,StatsRepo,Store

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/samber/lo"
)

type (
	TimeGenerator interface {
		Now() time.Time
	}

	StatsRepo interface {
		GetSensorsByHardware(ctx context.Context) (map[Hardware][]Sensor, error)
		GetCurrentSensorValues(ctx context.Context) (map[Sensor]float32, error)
	}

	Store interface {
		StoreValue(sID SensorID, value Value) error
		GetValuesForRange(sID SensorID, from, to time.Time) ([]Value, error)
	}
)

type Service struct {
	timeGenerator TimeGenerator
	statsRepo     StatsRepo
	store         Store
}

func NewService(timeGenerator TimeGenerator, statsRepo StatsRepo, store Store) (*Service, error) {
	if lo.IsNil(timeGenerator) {
		return nil, errors.New("time generator is nil")
	}
	if lo.IsNil(statsRepo) {
		return nil, errors.New("stats repo is nil")
	}
	if lo.IsNil(store) {
		return nil, errors.New("store is nil")
	}
	return &Service{
		timeGenerator: timeGenerator,
		statsRepo:     statsRepo,
		store:         store,
	}, nil
}

func (s *Service) GetStats(ctx context.Context, req GetStatsRequest) (GetStatsResponse, error) {
	zero := GetStatsResponse{}

	if err := req.Validate(); err != nil {
		return zero, fmt.Errorf("validate request: %w", err)
	}

	now := s.timeGenerator.Now()

	sensorsByHardware, err := s.statsRepo.GetSensorsByHardware(ctx)
	if err != nil {
		return zero, fmt.Errorf("get sensors: %w", err)
	}

	currentValues, err := s.statsRepo.GetCurrentSensorValues(ctx)
	if err != nil {
		return zero, fmt.Errorf("get current values: %w", err)
	}

	if err = s.storeValues(now, sensorsByHardware, currentValues); err != nil {
		return zero, fmt.Errorf("store values: %w", err)
	}

	return s.getValuesForRange(now, req, sensorsByHardware)
}

func (s *Service) storeValues(
	now time.Time,
	sensorsByHardware map[Hardware][]Sensor,
	currentValues map[Sensor]float32,
) error {
	for _, sensors := range sensorsByHardware {
		if len(sensors) == 0 {
			continue
		}

		for _, sensor := range sensors {
			currentSensorValue, ok := currentValues[sensor]
			if !ok {
				continue
			}

			value := Value{
				Value:     int64(math.Round(float64(currentSensorValue))),
				Timestamp: now,
			}
			if err := s.store.StoreValue(sensor.ID, value); err != nil {
				return fmt.Errorf("store value: %w", err)
			}
		}
	}

	return nil
}

func (s *Service) getValuesForRange(now time.Time, req GetStatsRequest, sensorsByHardware map[Hardware][]Sensor) (GetStatsResponse, error) {
	resp := GetStatsResponse{
		Stats: map[Hardware]map[SensorType]map[Sensor][]Value{},
	}

	for hardware, sensors := range sensorsByHardware {
		resp.Stats[hardware] = make(map[SensorType]map[Sensor][]Value, len(sensors))

		for _, sensor := range sensors {
			if _, ok := resp.Stats[hardware][sensor.Type]; !ok {
				resp.Stats[hardware][sensor.Type] = make(map[Sensor][]Value, len(sensors))
			}
			values, err := s.store.GetValuesForRange(sensor.ID, now.Add(-req.ForRange), now)
			if err != nil {
				return GetStatsResponse{}, fmt.Errorf("get values for range: %w", err)
			}
			resp.Stats[hardware][sensor.Type][sensor] = values
		}
	}

	return resp, nil
}
