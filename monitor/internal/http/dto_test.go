package http

import (
	"testing"
	"time"

	"github.com/genvmoroz/win-stats-service/internal/core"
	"github.com/stretchr/testify/require"
)

func TestRespFromCore(t *testing.T) {
	t.Parallel()

	var (
		cpu  = core.Hardware{ID: "/cpu0", Name: "INTEL CORE I7-7700K", Type: core.CPU}
		gpu  = core.Hardware{ID: "/gpu0", Name: "NVIDIA GEFORCE GTX 1080", Type: core.GPU}
		ram0 = core.Hardware{ID: "/ram0", Name: "KINGSTON DDR4 16GB", Type: core.RAM}
		ram1 = core.Hardware{ID: "/ram1", Name: "KINGSTON DDR4 16GB", Type: core.RAM}
	)

	var (
		cpuTemp0 = core.Sensor{ID: "/cpu0/core0/temperature", Name: "CPU Core 0 Temperature", Type: core.Temperature, MaxValue: 100}
		cpuTemp1 = core.Sensor{ID: "/cpu0/core1/temperature", Name: "CPU Core 1 Temperature", Type: core.Temperature, MaxValue: 100}

		cpuLoad0 = core.Sensor{ID: "/cpu0/core0/load", Name: "CPU Core 0 Load", Type: core.Load, MaxValue: 100}
		cpuLoad1 = core.Sensor{ID: "/cpu0/core1/load", Name: "CPU Core 1 Load", Type: core.Load, MaxValue: 100}

		gpuTemp0 = core.Sensor{ID: "/gpu0/core0/temperature", Name: "GPU Core 0 Temperature", Type: core.Temperature, MaxValue: 100}
		gpuTemp1 = core.Sensor{ID: "/gpu0/core1/temperature", Name: "GPU Core 1 Temperature", Type: core.Temperature, MaxValue: 100}

		ram0Load = core.Sensor{ID: "/ram0/load", Name: "RAM 0 Load", Type: core.Load, MaxValue: 100}
		ram1Load = core.Sensor{ID: "/ram1/load", Name: "RAM 1 Load", Type: core.Load, MaxValue: 100}
	)

	now := time.Now()

	coreResp := core.GetStatsResponse{
		Stats: map[core.Hardware]map[core.SensorType]map[core.Sensor][]core.Value{
			cpu: {
				core.Temperature: {
					cpuTemp0: {
						{Value: 50, Timestamp: now.Add(-10 * time.Second)},
						{Value: 51, Timestamp: now.Add(-5 * time.Second)},
						{Value: 52, Timestamp: now},
					},
					cpuTemp1: {
						{Value: 60, Timestamp: now.Add(-10 * time.Second)},
						{Value: 61, Timestamp: now.Add(-5 * time.Second)},
						{Value: 62, Timestamp: now},
					},
				},
				core.Load: {
					cpuLoad0: {
						{Value: 30, Timestamp: now.Add(-10 * time.Second)},
						{Value: 31, Timestamp: now.Add(-5 * time.Second)},
						{Value: 32, Timestamp: now},
					},
					cpuLoad1: {
						{Value: 40, Timestamp: now.Add(-10 * time.Second)},
						{Value: 41, Timestamp: now.Add(-5 * time.Second)},
						{Value: 42, Timestamp: now},
					},
				},
			},
			gpu: {
				core.Temperature: {
					gpuTemp0: {
						{Value: 80, Timestamp: now.Add(-10 * time.Second)},
						{Value: 81, Timestamp: now.Add(-5 * time.Second)},
						{Value: 82, Timestamp: now},
					},
					gpuTemp1: {
						{Value: 90, Timestamp: now.Add(-10 * time.Second)},
						{Value: 91, Timestamp: now.Add(-5 * time.Second)},
						{Value: 92, Timestamp: now},
					},
				},
			},
			ram0: {
				core.Load: {
					ram0Load: {
						{Value: 70, Timestamp: now.Add(-10 * time.Second)},
						{Value: 71, Timestamp: now.Add(-5 * time.Second)},
						{Value: 72, Timestamp: now},
					},
				},
			},
			ram1: {
				core.Load: {
					ram1Load: {
						{Value: 80, Timestamp: now.Add(-10 * time.Second)},
						{Value: 81, Timestamp: now.Add(-5 * time.Second)},
						{Value: 82, Timestamp: now},
					},
				},
			},
		},
	}

	want := GetStatsResponse{
		Stats: Stats{
			Hardware: []Hardware{
				{
					Name: "CPU: INTEL CORE I7-7700K [/cpu0]",
					SensorTypes: []SensorType{
						{
							TypeName: core.Load.String(),
							Unit:     core.Load.Unit().String(),
							Sensors: []Sensor{
								{
									Name:     "CPU Core 0 Load [/cpu0/core0/load]",
									MaxValue: 100,
									Values: []Value{
										{Value: 30, Timestamp: now.Add(-10 * time.Second).UnixMilli()},
										{Value: 31, Timestamp: now.Add(-5 * time.Second).UnixMilli()},
										{Value: 32, Timestamp: now.UnixMilli()},
									},
								},
								{
									Name:     "CPU Core 1 Load [/cpu0/core1/load]",
									MaxValue: 100,
									Values: []Value{
										{Value: 40, Timestamp: now.Add(-10 * time.Second).UnixMilli()},
										{Value: 41, Timestamp: now.Add(-5 * time.Second).UnixMilli()},
										{Value: 42, Timestamp: now.UnixMilli()},
									},
								},
							},
						},
						{
							TypeName: core.Temperature.String(),
							Unit:     core.Temperature.Unit().String(),
							Sensors: []Sensor{
								{
									Name:     "CPU Core 0 Temperature [/cpu0/core0/temperature]",
									MaxValue: 100,
									Values: []Value{
										{Value: 50, Timestamp: now.Add(-10 * time.Second).UnixMilli()},
										{Value: 51, Timestamp: now.Add(-5 * time.Second).UnixMilli()},
										{Value: 52, Timestamp: now.UnixMilli()},
									},
								},
								{
									Name:     "CPU Core 1 Temperature [/cpu0/core1/temperature]",
									MaxValue: 100,
									Values: []Value{
										{Value: 60, Timestamp: now.Add(-10 * time.Second).UnixMilli()},
										{Value: 61, Timestamp: now.Add(-5 * time.Second).UnixMilli()},
										{Value: 62, Timestamp: now.UnixMilli()},
									},
								},
							},
						},
					},
				},
				{
					Name: "GPU: NVIDIA GEFORCE GTX 1080 [/gpu0]",
					SensorTypes: []SensorType{
						{
							TypeName: core.Temperature.String(),
							Unit:     core.Temperature.Unit().String(),
							Sensors: []Sensor{
								{
									Name:     "GPU Core 0 Temperature [/gpu0/core0/temperature]",
									MaxValue: 100,
									Values: []Value{
										{Value: 80, Timestamp: now.Add(-10 * time.Second).UnixMilli()},
										{Value: 81, Timestamp: now.Add(-5 * time.Second).UnixMilli()},
										{Value: 82, Timestamp: now.UnixMilli()},
									},
								},
								{
									Name:     "GPU Core 1 Temperature [/gpu0/core1/temperature]",
									MaxValue: 100,
									Values: []Value{
										{Value: 90, Timestamp: now.Add(-10 * time.Second).UnixMilli()},
										{Value: 91, Timestamp: now.Add(-5 * time.Second).UnixMilli()},
										{Value: 92, Timestamp: now.UnixMilli()},
									},
								},
							},
						},
					},
				},
				{
					Name: "RAM: KINGSTON DDR4 16GB [/ram0]",
					SensorTypes: []SensorType{
						{
							TypeName: core.Load.String(),
							Unit:     core.Load.Unit().String(),
							Sensors: []Sensor{
								{
									Name:     "RAM 0 Load [/ram0/load]",
									MaxValue: 100,
									Values: []Value{
										{Value: 70, Timestamp: now.Add(-10 * time.Second).UnixMilli()},
										{Value: 71, Timestamp: now.Add(-5 * time.Second).UnixMilli()},
										{Value: 72, Timestamp: now.UnixMilli()},
									},
								},
							},
						},
					},
				},
				{
					Name: "RAM: KINGSTON DDR4 16GB [/ram1]",
					SensorTypes: []SensorType{
						{
							TypeName: core.Load.String(),
							Unit:     core.Load.Unit().String(),
							Sensors: []Sensor{
								{
									Name:     "RAM 1 Load [/ram1/load]",
									MaxValue: 100,
									Values: []Value{
										{Value: 80, Timestamp: now.Add(-10 * time.Second).UnixMilli()},
										{Value: 81, Timestamp: now.Add(-5 * time.Second).UnixMilli()},
										{Value: 82, Timestamp: now.UnixMilli()},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	require.Equal(t, want, fromCoreResp(coreResp))
}

func TestReqToCore(t *testing.T) {
	t.Parallel()

	type (
		input struct {
			in GetStatsRequest
		}
		want struct {
			out        core.GetStatsRequest
			errPresent bool
		}
	)
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "success with 1h",
			input: input{
				in: GetStatsRequest{
					Range: "1h",
				},
			},
			want: want{
				out: core.GetStatsRequest{
					ForRange: time.Hour,
				},
			},
		},
		{
			name: "success with 1m",
			input: input{
				in: GetStatsRequest{
					Range: "1m",
				},
			},
			want: want{
				out: core.GetStatsRequest{
					ForRange: time.Minute,
				},
			},
		},
		{
			name: "parse error",
			input: input{
				in: GetStatsRequest{
					Range: "1x-S",
				},
			},
			want: want{
				errPresent: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := toCoreGetStatsRequest(tt.input.in)
			require.Equal(t, tt.want.errPresent, err != nil)
			require.Equal(t, got, tt.want.out)
		})
	}
}
