//go:build windows

package ohm

import (
	"context"
	"fmt"

	"github.com/yusufpapurcu/wmi"
)

const (
	namespace           = "root\\LibreHardwareMonitor"
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
	return execQueryWithContext(ctx, r.getHardware(opts...))
}

func (r *Repo) GetSensors(ctx context.Context, opts ...SensorFilter) ([]Sensor, error) {
	return execQueryWithContext(ctx, r.getSensors(opts...))
}

func (r *Repo) getHardware(opts ...HardwareFilter) func() ([]Hardware, error) {
	return func() ([]Hardware, error) {
		query := selectHardwareQuery
		condition := narrowHardwareQuery(opts...)
		if len(condition) != 0 {
			query = fmt.Sprintf("%s %s", query, condition)
		}

		var res []Hardware
		if err := r.queryExecutor(query, &res, namespace); err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (r *Repo) getSensors(opts ...SensorFilter) func() ([]Sensor, error) {
	return func() ([]Sensor, error) {
		query := selectSensorsQuery
		condition := narrowSensorsQuery(opts...)
		if len(condition) != 0 {
			query = fmt.Sprintf("%s %s", query, condition)
		}

		var res []Sensor
		if err := r.queryExecutor(query, &res, namespace); err != nil {
			return nil, err
		}
		return res, nil
	}
}

func execQueryWithContext[T any](ctx context.Context, execFunc func() (T, error)) (T, error) {
	var zero T

	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	default:
	}

	type result struct {
		res T
		err error
	}

	resultChan := make(chan result, 1)

	go func() {
		res, err := execFunc()
		resultChan <- result{res, err}
	}()

	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case out := <-resultChan:
		return out.res, out.err
	}
}
