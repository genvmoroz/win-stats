package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	hardwareNameLabel = "hardwareName"
	hardwareTypeLabel = "hardwareType"
	sensorNameLabel   = "sensorName"
	sensorTypeLabel   = "sensorType"
)

type StatsReporter struct {
	sensorValueGaugeVec *prometheus.GaugeVec
}

func NewStatsReporter() *StatsReporter {
	return &StatsReporter{
		sensorValueGaugeVec: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "sensor_value",
				Help: "Sensor value",
			},
			[]string{hardwareNameLabel, hardwareTypeLabel, sensorNameLabel, sensorTypeLabel},
		),
	}
}

func (r *StatsReporter) Register() error {
	if err := prometheus.Register(r.sensorValueGaugeVec); err != nil {
		return fmt.Errorf("register sensor value gauge vec: %w", err)
	}

	return nil
}

func (r *StatsReporter) ReportSensorValue(value int64, hardwareName, hardwareType, sensorName, sensorType string) {
	r.sensorValueGaugeVec.
		With(
			map[string]string{
				hardwareNameLabel: hardwareName,
				hardwareTypeLabel: hardwareType,
				sensorNameLabel:   sensorName,
				sensorTypeLabel:   sensorType,
			},
		).
		Set(float64(value))
}
