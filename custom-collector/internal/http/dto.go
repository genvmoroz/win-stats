package http

import (
	"fmt"
	"sort"
	"time"

	"github.com/genvmoroz/custom-collector/internal/core"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type (
	GetStatsRequest struct {
		Range string `query:"range"`
	}

	GetStatsResponse struct {
		Stats Stats `json:"stats"`
	}

	Stats struct {
		Hardware []Hardware `json:"hardware,omitempty"`
	}

	Hardware struct {
		Name        string
		SensorTypes []SensorType `json:"sensorTypes,omitempty"`
	}

	SensorType struct {
		TypeName string
		Unit     string
		Sensors  []Sensor
	}

	Sensor struct {
		Name     string
		MaxValue int64
		Values   []Value
	}

	Value struct {
		Value     int64 `json:"value"`
		Timestamp int64 `json:"timestamp"`
	}
)

func fromCoreResp(in core.GetStatsResponse) GetStatsResponse {
	out := GetStatsResponse{}

	var (
		hwLen   = len(in.Stats)
		hwIndex = 0
	)
	out.Stats.Hardware = make([]Hardware, hwLen)
	for hw, sTypes := range in.Stats {
		out.Stats.Hardware[hwIndex] = Hardware{
			Name: hardwareName(hw),
		}

		var (
			sTypeLen   = len(sTypes)
			sTypeIndex = 0
		)
		out.Stats.Hardware[hwIndex].SensorTypes = make([]SensorType, sTypeLen)
		for sType, sensors := range sTypes {
			out.Stats.Hardware[hwIndex].SensorTypes[sTypeIndex] = SensorType{
				TypeName: sType.String(),
				Unit:     sType.Unit().String(),
			}

			var (
				sensorLen   = len(sensors)
				sensorIndex = 0
			)
			out.Stats.Hardware[hwIndex].SensorTypes[sTypeIndex].Sensors = make([]Sensor, sensorLen)
			for sensor, values := range sensors {
				out.Stats.Hardware[hwIndex].SensorTypes[sTypeIndex].Sensors[sensorIndex] = Sensor{
					Name:     sensorName(sensor),
					MaxValue: sensor.MaxValue,
					Values: lo.Map(values, func(v core.Value, _ int) Value {
						return Value{
							Value:     v.Value,
							Timestamp: v.Timestamp.UnixMilli(),
						}
					}),
				}
				sensorIndex++
			}
			sort.Slice(out.Stats.Hardware[hwIndex].SensorTypes[sTypeIndex].Sensors,
				func(i, j int) bool {
					return out.Stats.Hardware[hwIndex].SensorTypes[sTypeIndex].Sensors[i].Name <
						out.Stats.Hardware[hwIndex].SensorTypes[sTypeIndex].Sensors[j].Name
				},
			)
			sTypeIndex++
		}
		sort.Slice(out.Stats.Hardware[hwIndex].SensorTypes,
			func(i, j int) bool {
				return out.Stats.Hardware[hwIndex].SensorTypes[i].TypeName <
					out.Stats.Hardware[hwIndex].SensorTypes[j].TypeName
			},
		)
		hwIndex++
	}
	sort.Slice(out.Stats.Hardware,
		func(i, j int) bool {
			return out.Stats.Hardware[i].Name < out.Stats.Hardware[j].Name
		},
	)

	return out
}

func toCoreGetStatsRequest(in GetStatsRequest) (core.GetStatsRequest, error) {
	var zero core.GetStatsRequest

	duration, err := time.ParseDuration(in.Range)
	if err != nil {
		return zero, fmt.Errorf("parse duration: %w", err)
	}

	return core.GetStatsRequest{
		ForRange: duration,
	}, nil
}

func hardwareName(hw core.Hardware) string {
	return fmt.Sprintf("%s: %s [%s]", hw.Type, hw.Name, hw.ID)
}

func sensorName(s core.Sensor) string {
	return fmt.Sprintf("%s [%s]", s.Name, s.ID)
}

func parseGetStatsRequest(c echo.Context) (GetStatsRequest, error) {
	var req GetStatsRequest
	if err := c.Bind(&req); err != nil {
		return GetStatsRequest{}, err
	}

	return req, nil
}
