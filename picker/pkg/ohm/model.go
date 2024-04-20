package ohm

type Hardware struct {
	Name         string
	Identifier   string
	HardwareType HardwareType
	Parent       string
}

type Sensor struct {
	Name       string
	Identifier string
	SensorType SensorType
	Parent     string
	Value      float32
	Min        float32
	Max        float32
	Index      int
}
