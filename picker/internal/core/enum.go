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
	Battery
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
	Energy
	Current
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
