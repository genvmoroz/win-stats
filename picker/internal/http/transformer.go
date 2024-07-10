package http

import (
	"github.com/genvmoroz/win-stats-picker/internal/core"
	openapi "github.com/genvmoroz/win-stats-picker/internal/http/generated"
	"github.com/samber/lo"
)

type Transformer struct{}

func (t Transformer) GetStatsResponseFromCore(in core.GetStatsResponse) openapi.Stats {
	out := openapi.Stats{}
	if in.Stats == nil {
		return out
	}

	hardware := make([]openapi.Hardware, 0, len(in.Stats))
	for coreHW, coreSensors := range in.Stats {
		if len(coreSensors) == 0 {
			continue
		}

		hw := t.hardwareFromCore(coreHW)

		sensors := make([]openapi.Sensor, len(coreSensors))
		for idx, coreSensor := range coreSensors {
			sensors[idx] = t.sensorFromCore(coreSensor)
		}
		hw.Sensors = &sensors

		hardware = append(hardware, hw)
	}

	out.Hardware = &hardware

	return out
}

func (t Transformer) hardwareFromCore(in core.Hardware) openapi.Hardware {
	return openapi.Hardware{
		ID:   lo.ToPtr(string(in.ID)),
		Name: lo.ToPtr(in.Name),
		Type: lo.ToPtr(in.Type.String()),
	}
}

func (t Transformer) sensorFromCore(in core.Sensor) openapi.Sensor {
	return openapi.Sensor{
		ID:    lo.ToPtr(string(in.ID)),
		Name:  lo.ToPtr(in.Name),
		Type:  lo.ToPtr(in.Type.String()),
		Value: lo.ToPtr(t.valueFromCore(in.Value)),
	}
}

func (t Transformer) valueFromCore(in core.SensorValue) openapi.SensorValue {
	return openapi.SensorValue{
		Value:     lo.ToPtr(in.Value),
		Timestamp: lo.ToPtr(in.Timestamp.Unix()),
	}
}
