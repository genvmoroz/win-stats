package stats

import (
	"errors"
	"fmt"
	"math"

	"github.com/genvmoroz/win-stats/picker/internal/core"
	"github.com/genvmoroz/win-stats/picker/pkg/ohm"
)

func (r *Repo) toCoreHardware(in []ohm.Hardware) ([]core.Hardware, error) {
	out := make([]core.Hardware, 0, len(in))
	var errs []error
	for _, h := range in {
		t, err := toCoreHardwareType(h.HardwareType)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		out = append(out,
			core.Hardware{
				ID:   core.HardwareID(h.Identifier),
				Name: h.Name,
				Type: t,
			},
		)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return out, nil
}

func (r *Repo) toCoreSensors(in []ohm.Sensor) ([]core.Sensor, error) {
	out := make([]core.Sensor, 0, len(in))
	var errs []error

	for _, s := range in {
		transformed, err := r.toCoreSensor(s)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		out = append(out, transformed)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return out, nil
}

func (r *Repo) toCoreSensor(in ohm.Sensor) (core.Sensor, error) {
	t, err := toCoreSensorType(in.SensorType)
	if err != nil {
		return core.Sensor{}, err
	}
	return core.Sensor{
		ID:         core.SensorID(in.Identifier),
		HardwareID: core.HardwareID(in.Parent),
		Name:       in.Name,
		Type:       t,
		Value: core.SensorValue{
			Value:     int64(math.Round(float64(in.Value))),
			Timestamp: r.timegen.Now(),
		},
	}, nil
}

func toCoreHardwareType(in ohm.HardwareType) (core.HardwareType, error) {
	switch in {
	case ohm.Mainboard, ohm.Motherboard:
		return core.Motherboard, nil
	case ohm.SuperIO:
		return core.SuperIO, nil
	case ohm.CPU:
		return core.CPU, nil
	case ohm.GpuNvidia, ohm.GpuAti, ohm.GpuAmd, ohm.GpuIntel:
		return core.GPU, nil
	case ohm.TBalancer:
		return core.TBalancer, nil
	case ohm.HeatMaster:
		return core.HeatMaster, nil
	case ohm.HDD:
		return core.HDD, nil
	case ohm.RAM:
		return core.RAM, nil
	case ohm.Network:
		return core.Network, nil
	case ohm.Memory:
		return core.Memory, nil
	case ohm.Storage:
		return core.Storage, nil
	case ohm.Battery:
		return core.Battery, nil
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
	case ohm.Factor:
		return core.Factor, nil
	case ohm.Energy:
		return core.Energy, nil
	case ohm.Current:
		return core.Current, nil
	default:
		return core.UnknownSensorType, fmt.Errorf("unknown sensor type: %s", in)
	}
}
