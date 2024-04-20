package core

type Stats struct {
	Hardware []Hardware
}

type Hardware struct {
	ID      string
	Name    string
	Sensors []Sensor
	Type    string
}

type Sensor struct {
	ID   string
	Name string
	Type string

	Value SensorValue
}

type SensorValue struct {
	Timestamp int64
	Value     int64
}
