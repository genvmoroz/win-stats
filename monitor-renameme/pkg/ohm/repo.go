//go:build windows

package ohm

import (
	"context"
	"fmt"
)

const (
	namespace           = "root\\OpenHardwareMonitor"
	selectHardwareQuery = "SELECT * FROM Hardware"
	selectSensorsQuery  = "SELECT * FROM Sensor"
)

type Repo struct {
	queryExecutor func(query string, dst any, namespace string) error
}

func NewRepo() *Repo {
	return &Repo{
		queryExecutor: wmi.QueryNamespace,
	}
}

func (r *Repo) GetHardware(ctx context.Context, opts ...HardwareFilter) ([]Hardware, error) {
	return execQueryWithContext[[]Hardware](ctx, r.getHardware(opts...))
}

func (r *Repo) GetSensors(ctx context.Context, opts ...SensorFilter) ([]Sensor, error) {
	return execQueryWithContext[[]Sensor](ctx, r.getSensors(opts...))
}

func (r *Repo) getHardware(opts ...HardwareFilter) func(resChan chan []Hardware, errChan chan error) {
	return func(resChan chan []Hardware, errChan chan error) {
		query := selectHardwareQuery
		condition := narrowHardwareQuery(opts...)
		if len(condition) != 0 {
			query = fmt.Sprintf("%s %s", query, condition)
		}

		var res []Hardware
		if err := r.queryExecutor(query, &res, namespace); err != nil {
			errChan <- fmt.Errorf("exec query: %w", err)
		} else {
			resChan <- res
		}
	}
}

func (r *Repo) getSensors(opts ...SensorFilter) func(resChan chan []Sensor, errChan chan error) {
	return func(resChan chan []Sensor, errChan chan error) {
		query := selectSensorsQuery
		condition := narrowSensorsQuery(opts...)
		if len(condition) != 0 {
			query = fmt.Sprintf("%s %s", query, condition)
		}

		var res []Sensor
		if err := r.queryExecutor(query, &res, namespace); err != nil {
			errChan <- fmt.Errorf("exec query: %w", err)
		} else {
			resChan <- res
		}
	}
}

func execQueryWithContext[T any](ctx context.Context, execFunc func(resChan chan T, errChan chan error)) (T, error) {
	var zero T

	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	default:
	}

	resChan := make(chan T, 1)
	errChan := make(chan error, 1)

	go execFunc(resChan, errChan)

	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case res := <-resChan:
		return res, nil
	case err := <-errChan:
		return zero, err
	}
}
