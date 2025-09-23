package stats

import (
	"context"
	"errors"

	"github.com/genvmoroz/win-stats/picker/internal/core"
	"github.com/samber/lo"
	"golang.org/x/sync/singleflight"
)

type SingleflightRepo struct {
	baseRepo core.StatsRepo
	group    singleflight.Group
}

func NewSingleflightRepo(baseRepo core.StatsRepo) (*SingleflightRepo, error) {
	if lo.IsNil(baseRepo) {
		return nil, errors.New("base repo is nil")
	}

	return &SingleflightRepo{
		baseRepo: baseRepo,
		group:    singleflight.Group{},
	}, nil
}

const getSensorsByHardwareKey = "GetSensorsByHardware"

func (c *SingleflightRepo) GetSensorsByHardware(ctx context.Context) (map[core.Hardware][]core.Sensor, error) {
	result, err, _ := c.group.Do(getSensorsByHardwareKey, func() (any, error) {
		var (
			stats map[core.Hardware][]core.Sensor
			err   error
		)
		stats, err = c.baseRepo.GetSensorsByHardware(ctx)
		return stats, err
	})
	if err != nil {
		return nil, err
	}
	stats, ok := result.(map[core.Hardware][]core.Sensor)
	if !ok {
		return nil, errors.New("invalid result type, must never happen, check the code")
	}
	return stats, nil
}
