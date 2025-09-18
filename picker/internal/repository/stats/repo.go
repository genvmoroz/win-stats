package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/genvmoroz/win-stats/picker/internal/core"
	"github.com/genvmoroz/win-stats/picker/pkg/ohm"
	"github.com/samber/lo"
)

type TimeGenerator interface {
	Now() time.Time
}

type Repo struct {
	ohm     *ohm.Repo
	timegen TimeGenerator
}

func NewRepo(timegen TimeGenerator) (*Repo, error) {
	if lo.IsNil(timegen) {
		return nil, fmt.Errorf("time generator is nil")
	}
	return &Repo{
		ohm:     ohm.NewRepo(),
		timegen: timegen,
	}, nil
}

func (r *Repo) GetSensorsByHardware(ctx context.Context) (map[core.Hardware][]core.Sensor, error) {
	hardware, err := r.getHardware(ctx)
	if err != nil {
		return nil, fmt.Errorf("get hardware: %w", err)
	}

	sensors, err := r.ohm.GetSensors(ctx)
	if err != nil {
		return nil, fmt.Errorf("get sensors: %w", err)
	}

	transformedSensors, err := r.toCoreSensors(sensors)
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

func (r *Repo) getHardware(ctx context.Context) ([]core.Hardware, error) {
	hardware, err := r.ohm.GetHardware(ctx)
	if err != nil {
		return nil, fmt.Errorf("get hardware: %w", err)
	}

	transformed, err := r.toCoreHardware(hardware)
	if err != nil {
		return nil, fmt.Errorf("transform hardware: %w", err)
	}

	return transformed, nil
}
