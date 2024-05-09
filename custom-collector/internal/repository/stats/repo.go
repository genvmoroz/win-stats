package stats

import (
	"context"

	"github.com/genvmoroz/custom-collector/internal/core"
)

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) GetHardware(ctx context.Context) ([]core.Hardware, error) {
	return nil, nil
}

// todo: add ability to fetch sensors indeed for provided hardware IDs
func (r *Repo) GetSensorsByHardware(ctx context.Context) (map[core.Hardware][]core.Sensor, error) {
	return nil, nil
}

func (r *Repo) GetCurrentSensorValues(ctx context.Context) (map[core.Sensor]float32, error) {
	return nil, nil
}
