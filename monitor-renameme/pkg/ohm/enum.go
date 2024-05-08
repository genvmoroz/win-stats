package ohm

type HardwareType string

const (
	Mainboard  HardwareType = "Mainboard"
	SuperIO    HardwareType = "SuperIO"
	CPU        HardwareType = "CPU"
	GpuNvidia  HardwareType = "GpuNvidia"
	GpuAti     HardwareType = "GpuAti"
	TBalancer  HardwareType = "TBalancer"
	HeatMaster HardwareType = "HeatMaster"
	HDD        HardwareType = "HDD"
	RAM        HardwareType = "RAM"
)

type SensorType string

const (
	Voltage     SensorType = "Voltage"     // Volt
	Clock       SensorType = "Clock"       // Megahertz
	Temperature SensorType = "Temperature" // Celsius
	Load        SensorType = "Load"        // Percentage
	Fan         SensorType = "Fan"         // Revolutions per minute
	Flow        SensorType = "Flow"        // Liters per hour
	Control     SensorType = "Control"     // Percentage
	Level       SensorType = "Level"
	Power       SensorType = "Power"      // ?
	SmallData   SensorType = "SmallData"  // ?
	Throughput  SensorType = "Throughput" // ?
	Data        SensorType = "Data"       // ?
)
