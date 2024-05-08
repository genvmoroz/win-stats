package stats

import (
	"fmt"
	"math"

	"github.com/genvmoroz/win-stats-service/internal/core"
	"github.com/genvmoroz/win-stats-service/pkg/ohm"
)

func toCoreHardware(in []ohm.Hardware) ([]core.Hardware, error) {
	out := make([]core.Hardware, len(in))
	for i, h := range in {
		t, err := toCoreHardwareType(h.HardwareType)
		if err != nil {
			return nil, fmt.Errorf("transform hardware type: %w", err)
		}
		out[i] = core.Hardware{
			ID:   core.HardwareID(h.Identifier),
			Name: h.Name,
			Type: t,
		}
	}
	return out, nil
}

func toCoreSensors(in []ohm.Sensor) ([]core.Sensor, error) {
	out := make([]core.Sensor, len(in))
	for i, s := range in {
		transformed, err := toCoreSensor(s)
		if err != nil {
			return nil, fmt.Errorf("transform sensor: %w", err)
		}
		out[i] = transformed
	}
	return out, nil
}

func toCoreSensor(in ohm.Sensor) (core.Sensor, error) {
	t, err := toCoreSensorType(in.SensorType)
	if err != nil {
		return core.Sensor{}, fmt.Errorf("transform sensor type: %w", err)
	}
	return core.Sensor{
		ID:         core.SensorID(in.Identifier),
		HardwareID: core.HardwareID(in.Parent),
		Name:       in.Name,
		Type:       t,
		MaxValue:   int64(math.Round(float64(in.Max))),
	}, nil
}

func toCoreHardwareType(in ohm.HardwareType) (core.HardwareType, error) {
	switch in {
	case ohm.Mainboard:
		return core.Motherboard, nil
	case ohm.SuperIO:
		return core.SuperIO, nil
	case ohm.CPU:
		return core.CPU, nil
	case ohm.GpuNvidia, ohm.GpuAti:
		return core.GPU, nil
	case ohm.TBalancer:
		return core.TBalancer, nil
	case ohm.HeatMaster:
		return core.HeatMaster, nil
	case ohm.HDD:
		return core.HDD, nil
	case ohm.RAM:
		return core.RAM, nil
	default:
		return core.UnknownHardwareType, fmt.Errorf("unknown hardware type: %s", in)
	}
}

func toCoreSensorType(in ohm.SensorType) (core.SensorType, error) {
	switch in {
	case ohm.Voltage:
		return core.Voltage, nil
	case ohm.Clock:
		return core.Clock, nil
	case ohm.Temperature:
		return core.Temperature, nil
	case ohm.Load:
		return core.Load, nil
	case ohm.Fan:
		return core.Fan, nil
	case ohm.Flow:
		return core.Flow, nil
	case ohm.Control:
		return core.Control, nil
	case ohm.Level:
		return core.Level, nil
	case ohm.Power:
		return core.Power, nil
	case ohm.SmallData:
		return core.SmallData, nil
	case ohm.Throughput:
		return core.Throughput, nil
	case ohm.Data:
		return core.Data, nil
	default:
		return core.UnknownSensorType, fmt.Errorf("unknown sensor type: %s", in)
	}
}
