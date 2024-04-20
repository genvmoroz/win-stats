package picker

import (
	"fmt"

	"github.com/genvmoroz/win-stats-prometheus-collector/internal/core"
	openapi "github.com/genvmoroz/win-stats-prometheus-collector/internal/repository/picker/generated"
	"github.com/samber/lo"
)

type Transformer struct{}

func (t Transformer) GetStatsResponseFromOpenAPI(resp openapi.GetStatsResponse) (core.Stats, error) {
	if resp.JSON200 == nil {
		return core.Stats{}, fmt.Errorf("response body is nil")
	}
	return t.statsFromOpenAPI(*resp.JSON200), nil
}

func (t Transformer) statsFromOpenAPI(in openapi.Stats) core.Stats {
	return core.Stats{
		Hardware: t.multipleHardwareFromOpenAPI(lo.FromPtr(in.Hardware)),
	}
}

func (t Transformer) multipleHardwareFromOpenAPI(in []openapi.Hardware) []core.Hardware {
	if in == nil {
		return nil
	}

	out := make([]core.Hardware, len(in))
	for idx, hw := range in {
		out[idx] = t.hardwareFromOpenAPI(hw)
	}

	return out
}

func (t Transformer) hardwareFromOpenAPI(in openapi.Hardware) core.Hardware {
	return core.Hardware{
		ID:      lo.FromPtr(in.ID),
		Name:    lo.FromPtr(in.Name),
		Sensors: t.sensorsFromOpenAPI(lo.FromPtr(in.Sensors)),
		Type:    lo.FromPtr(in.Type),
	}
}

func (t Transformer) sensorsFromOpenAPI(in []openapi.Sensor) []core.Sensor {
	if in == nil {
		return nil
	}

	out := make([]core.Sensor, len(in))
	for idx, s := range in {
		out[idx] = t.sensorFromOpenAPI(s)
	}

	return out
}

func (t Transformer) sensorFromOpenAPI(in openapi.Sensor) core.Sensor {
	return core.Sensor{
		ID:    lo.FromPtr(in.ID),
		Name:  lo.FromPtr(in.Name),
		Type:  lo.FromPtr(in.Type),
		Value: t.sensorValueFromOpenAPI(lo.FromPtr(in.Value)),
	}
}

func (t Transformer) sensorValueFromOpenAPI(in openapi.SensorValue) core.SensorValue {
	return core.SensorValue{
		Timestamp: lo.FromPtr(in.Timestamp),
		Value:     lo.FromPtr(in.Value),
	}
}

func (t Transformer) HealthCheckFromOpenAPI(_ openapi.HealthCheckResponse) (struct{}, error) {
	return struct{}{}, nil
}
