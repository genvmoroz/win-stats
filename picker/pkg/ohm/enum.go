package ohm

type HardwareType string

const (
	Mainboard   HardwareType = "Mainboard"
	SuperIO     HardwareType = "SuperIO"
	CPU         HardwareType = "Cpu"
	GpuNvidia   HardwareType = "GpuNvidia"
	GpuAti      HardwareType = "GpuAti"
	GpuAmd      HardwareType = "GpuAmd"
	GpuIntel    HardwareType = "GpuIntel"
	TBalancer   HardwareType = "TBalancer"
	HeatMaster  HardwareType = "HeatMaster"
	HDD         HardwareType = "HDD"
	RAM         HardwareType = "RAM"
	Network     HardwareType = "Network"
	Memory      HardwareType = "Memory"
	Storage     HardwareType = "Storage"
	Motherboard HardwareType = "Motherboard"
	Battery     HardwareType = "Battery"
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
	Factor      SensorType = "Factor"
	Energy      SensorType = "Energy"
)
