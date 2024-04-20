package core

import (
	"time"
)

type GetStatsResponse struct {
	Stats map[Hardware][]Sensor
}

type (
	HardwareID string
	SensorID   string
)

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
	Value      SensorValue
}

type SensorValue struct {
	Value     int64
	Timestamp time.Time
}
