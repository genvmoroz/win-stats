// Package prometheus provides reporting to Prometheus metrics.
package prometheus

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	hostLabel         = "host"
	hardwareIDLabel   = "hardwareID"
	hardwareNameLabel = "hardwareName"
	hardwareTypeLabel = "hardwareType"
	sensorIDLabel     = "sensorID"
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
			[]string{hostLabel, hardwareIDLabel, hardwareNameLabel, hardwareTypeLabel, sensorIDLabel, sensorNameLabel, sensorTypeLabel},
		),
	}
}

func (r *StatsReporter) Register() error {
	if err := prometheus.Register(r.sensorValueGaugeVec); err != nil {
		return fmt.Errorf("register sensor value gauge vec: %w", err)
	}

	return nil
}

func (r *StatsReporter) ReportSensorValue(value int64, host, hardwareID, hardwareName, hardwareType, sensorID, sensorName, sensorType string) {
	r.sensorValueGaugeVec.
		With(
			map[string]string{
				hostLabel:         host,
				hardwareIDLabel:   hardwareID,
				hardwareNameLabel: hardwareName,
				hardwareTypeLabel: hardwareType,
				sensorIDLabel:     sensorID,
				sensorNameLabel:   sensorName,
				sensorTypeLabel:   sensorType,
			},
		).
		Set(float64(value))
}
