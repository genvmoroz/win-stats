package core

//go:generate stringer -output=enum_strings.go -type=HardwareType,SensorType,Unit

type HardwareType int

const (
	UnknownHardwareType HardwareType = iota
	SuperIO
	CPU
	GPU
	TBalancer
	HeatMaster
	HDD
	RAM
	Network
	Memory
	Storage
	Motherboard
)

type SensorType int

const (
	UnknownSensorType SensorType = iota
	Voltage
	Clock
	Temperature
	Load
	Fan
	Flow
	Control
	Level
	Power
	SmallData
	Throughput
	Data
	Factor
)

type Unit int

const (
	UnknownUnit Unit = iota
	Volt
	Megahertz
	Celsius
	Percentage
	RevolutionsPerMinute
	LitersPerHour
	Watts
	Gigabytes
	Megabytes
	KilobytesPerSecond
)

func (st SensorType) Unit() Unit {
	switch st {
	case UnknownSensorType:
		return UnknownUnit
	case Voltage:
		return Volt
	case Clock:
		return Megahertz
	case Temperature:
		return Celsius
	case Load:
		return Percentage
	case Fan:
		return RevolutionsPerMinute
	case Flow:
		return LitersPerHour
	case Control:
		return Percentage
	case Level:
		return Percentage
	case Power:
		return Watts
	case SmallData:
		return Megabytes
	case Throughput:
		return KilobytesPerSecond
	case Data:
		return Gigabytes
	default:
		return UnknownUnit
	}
}
