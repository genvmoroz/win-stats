package core

import (
	"fmt"
	"time"
)

type (
	GetStatsRequest struct {
		ForRange time.Duration
	}

	GetStatsResponse struct {
		Stats map[Hardware]map[SensorType]map[Sensor][]Value
	}
)

type (
	HardwareID string
	SensorID   string
)

type Value struct {
	Value     int64
	Timestamp time.Time
}

type Hardware struct {
	ID   HardwareID
	Name string
	Type HardwareType
}

type Sensor struct {
	ID         SensorID
	HardwareID HardwareID
	Name       string
	Type       SensorType
	MaxValue   int64
}

func (r GetStatsRequest) Validate() error {
	if r.ForRange <= 0 {
		return fmt.Errorf("range must be greater than 0")
	}
	return nil
}
