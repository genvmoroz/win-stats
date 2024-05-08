package stats

import (
	"context"
	"fmt"

	"github.com/genvmoroz/win-stats-service/internal/core"
	"github.com/genvmoroz/win-stats-service/pkg/ohm"
)

type Repo struct {
	ohm *ohm.Repo
}

func NewRepo() *Repo {
	return &Repo{
		ohm: ohm.NewRepo(),
	}
}

func (r *Repo) GetHardware(ctx context.Context) ([]core.Hardware, error) {
	hardware, err := r.ohm.GetHardware(ctx)
	if err != nil {
		return nil, fmt.Errorf("get hardware: %w", err)
	}

	transformed, err := toCoreHardware(hardware)
	if err != nil {
		return nil, fmt.Errorf("transform hardware: %w", err)
	}

	return transformed, nil
}

// todo: add ability to fetch sensors indeed for provided hardware IDs
func (r *Repo) GetSensorsByHardware(ctx context.Context) (map[core.Hardware][]core.Sensor, error) {
	hardware, err := r.GetHardware(ctx)
	if err != nil {
		return nil, fmt.Errorf("get hardware: %w", err)
	}

	sensors, err := r.ohm.GetSensors(ctx)
	if err != nil {
		return nil, fmt.Errorf("get sensors: %w", err)
	}

	transformedSensors, err := toCoreSensors(sensors)
	if err != nil {
		return nil, fmt.Errorf("transform sensors: %w", err)
	}

	result := make(map[core.Hardware][]core.Sensor, len(hardware))

	for _, hw := range hardware {
		for _, sensor := range transformedSensors {
			if sensor.HardwareID == hw.ID {
				result[hw] = append(result[hw], sensor)
			}
		}
	}

	return result, nil
}

func (r *Repo) GetCurrentSensorValues(ctx context.Context) (map[core.Sensor]float32, error) {
	sensors, err := r.ohm.GetSensors(ctx)
	if err != nil {
		return nil, fmt.Errorf("get sensors: %w", err)
	}

	result := make(map[core.Sensor]float32, len(sensors))
	for _, sensor := range sensors {
		transformed, err := toCoreSensor(sensor)
		if err != nil {
			return nil, fmt.Errorf("transform sensor: %w", err)
		}
		if _, ok := result[transformed]; ok {
			return nil, fmt.Errorf("duplicated sensor by ID (%+v): %w", transformed, err) // should never happen
		}
		result[transformed] = sensor.Value
	}

	return result, nil
}
