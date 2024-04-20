package core_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/genvmoroz/win-stats-service/internal/core"
	"github.com/genvmoroz/win-stats-service/internal/core/mock"
	"github.com/genvmoroz/win-stats-service/internal/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"math/rand/v2"
)

func TestServiceGetStats(t *testing.T) {
	t.Parallel()

	tests := []testutils.Test[*core.Service, testData, testDeps]{
		{
			Desc:     "success",
			EditData: nil,
			EditFlow: nil,
			TestFunc: func(service *core.Service, data testData) {
				got, err := service.GetStats(data.ctx, data.req)
				require.NoError(t, err)
				require.Equal(t, data.resp, got)
			},
		},
		// todo: add more test cases
	}
	for _, test := range tests {
		test.Run(t, initTestDeps, initTestData, initTestHookSet)
	}
}

type testData struct {
	ctx  context.Context
	req  core.GetStatsRequest
	resp core.GetStatsResponse

	now time.Time

	sensorsByHardware           map[core.Hardware][]core.Sensor
	currentValuesBySensor       map[core.Sensor]float32
	valuesPerSensorsForHardware map[core.Hardware]map[core.SensorType]map[core.Sensor][]core.Value
}

func initTestData(t *testing.T) testData {
	t.Helper()

	now := time.Now()

	var (
		cpu  = core.Hardware{ID: "/cpu0", Name: "INTEL CORE I7-7700K", Type: core.CPU}
		gpu  = core.Hardware{ID: "/gpu0", Name: "NVIDIA GEFORCE GTX 1080", Type: core.GPU}
		ram0 = core.Hardware{ID: "/ram0", Name: "KINGSTON DDR4 16GB", Type: core.RAM}
		ram1 = core.Hardware{ID: "/ram1", Name: "KINGSTON DDR4 16GB", Type: core.RAM}
	)

	var (
		cpu0Clock = core.Sensor{ID: "/cpu0/core0/clock", Name: "Core 0 Clock", Type: core.Clock, MaxValue: 100}
		cpu1Clock = core.Sensor{ID: "/cpu0/core1/clock", Name: "Core 1 Clock", Type: core.Clock, MaxValue: 100}
		cpu0Temp  = core.Sensor{ID: "/cpu0/core0/temp", Name: "Core 0 Temperature", Type: core.Temperature, MaxValue: 100}
		cpu1Temp  = core.Sensor{ID: "/cpu0/core1/temp", Name: "Core 1 Temperature", Type: core.Temperature, MaxValue: 100}
		gpu0Clock = core.Sensor{ID: "/gpu0/core0/clock", Name: "Core 0 Clock", Type: core.Clock, MaxValue: 100}
		gpu1Clock = core.Sensor{ID: "/gpu0/core1/clock", Name: "Core 1 Clock", Type: core.Clock, MaxValue: 100}
		gpu0Temp  = core.Sensor{ID: "/gpu0/core0/temp", Name: "Core 0 Temperature", Type: core.Temperature, MaxValue: 100}
		gpu1Temp  = core.Sensor{ID: "/gpu0/core1/temp", Name: "Core 1 Temperature", Type: core.Temperature, MaxValue: 100}
		ram0Usage = core.Sensor{ID: "/ram0/load", Name: "RAM 0 Load", Type: core.Load, MaxValue: 100}
		ram1Usage = core.Sensor{ID: "/ram1/load", Name: "RAM 1 Load", Type: core.Load, MaxValue: 100}
	)

	storedValues := func(current core.Value) []core.Value {
		return append(generateValuesForRange(now.Add(-time.Hour), now, time.Minute), current)
	}

	var (
		currentCPU0ClockValue = core.Value{Value: int64(850), Timestamp: now}
		currentCPU1ClockValue = core.Value{Value: int64(860), Timestamp: now}
		currentCPU0TempValue  = core.Value{Value: int64(42), Timestamp: now}
		currentCPU1TempValue  = core.Value{Value: int64(43), Timestamp: now}
		currentGPU0ClockValue = core.Value{Value: int64(1550), Timestamp: now}
		currentGPU1ClockValue = core.Value{Value: int64(1560), Timestamp: now}
		currentGPU0TempValue  = core.Value{Value: int64(52), Timestamp: now}
		currentGPU1TempValue  = core.Value{Value: int64(53), Timestamp: now}
		currentRAM0UsageValue = core.Value{Value: int64(40), Timestamp: now}
		currentRAM1UsageValue = core.Value{Value: int64(45), Timestamp: now}
	)

	var (
		clocksCPU0 = storedValues(core.Value{Value: currentCPU0ClockValue.Value, Timestamp: now})
		clocksCPU1 = storedValues(core.Value{Value: currentCPU1ClockValue.Value, Timestamp: now})
		tempsCPU0  = storedValues(core.Value{Value: currentCPU0TempValue.Value, Timestamp: now})
		tempsCPU1  = storedValues(core.Value{Value: currentCPU1TempValue.Value, Timestamp: now})
		clocksGPU0 = storedValues(core.Value{Value: currentGPU0ClockValue.Value, Timestamp: now})
		clocksGPU1 = storedValues(core.Value{Value: currentGPU1ClockValue.Value, Timestamp: now})
		tempsGPU0  = storedValues(core.Value{Value: currentGPU0TempValue.Value, Timestamp: now})
		tempsGPU1  = storedValues(core.Value{Value: currentGPU1TempValue.Value, Timestamp: now})
		loadsRAM0  = storedValues(core.Value{Value: currentRAM0UsageValue.Value, Timestamp: now})
		loadsRAM1  = storedValues(core.Value{Value: currentRAM1UsageValue.Value, Timestamp: now})
	)

	return testData{
		ctx: context.Background(),
		req: core.GetStatsRequest{
			ForRange: time.Hour,
		},
		resp: core.GetStatsResponse{
			Stats: map[core.Hardware]map[core.SensorType]map[core.Sensor][]core.Value{
				cpu: {
					core.Clock: {
						cpu0Clock: clocksCPU0,
						cpu1Clock: clocksCPU1,
					},
					core.Temperature: {
						cpu0Temp: tempsCPU0,
						cpu1Temp: tempsCPU1,
					},
				},
				gpu: {
					core.Clock: {
						gpu0Clock: clocksGPU0,
						gpu1Clock: clocksGPU1,
					},
					core.Temperature: {
						gpu0Temp: tempsGPU0,
						gpu1Temp: tempsGPU1,
					},
				},
				ram0: {
					core.Load: {ram0Usage: loadsRAM0},
				},
				ram1: {
					core.Load: {ram1Usage: loadsRAM1},
				},
			},
		},
		now: now,
		sensorsByHardware: map[core.Hardware][]core.Sensor{
			cpu:  {cpu0Clock, cpu1Clock, cpu0Temp, cpu1Temp},
			gpu:  {gpu0Clock, gpu1Clock, gpu0Temp, gpu1Temp},
			ram0: {ram0Usage},
			ram1: {ram1Usage},
		},
		currentValuesBySensor: map[core.Sensor]float32{
			cpu0Clock: float32(currentCPU0ClockValue.Value),
			cpu1Clock: float32(currentCPU1ClockValue.Value),
			cpu0Temp:  float32(currentCPU0TempValue.Value),
			cpu1Temp:  float32(currentCPU1TempValue.Value),
			gpu0Clock: float32(currentGPU0ClockValue.Value),
			gpu1Clock: float32(currentGPU1ClockValue.Value),
			gpu0Temp:  float32(currentGPU0TempValue.Value),
			gpu1Temp:  float32(currentGPU1TempValue.Value),
			ram0Usage: float32(currentRAM0UsageValue.Value),
			ram1Usage: float32(currentRAM1UsageValue.Value),
		},
		valuesPerSensorsForHardware: map[core.Hardware]map[core.SensorType]map[core.Sensor][]core.Value{
			cpu: {
				core.Clock: {
					cpu0Clock: clocksCPU0,
					cpu1Clock: clocksCPU1,
				},
				core.Temperature: {
					cpu0Temp: tempsCPU0,
					cpu1Temp: tempsCPU1,
				},
			},
			gpu: {
				core.Clock: {
					gpu0Clock: clocksGPU0,
					gpu1Clock: clocksGPU1,
				},
				core.Temperature: {
					gpu0Temp: tempsGPU0,
					gpu1Temp: tempsGPU1,
				},
			},
			ram0: {
				core.Load: {ram0Usage: loadsRAM0},
			},
			ram1: {
				core.Load: {ram1Usage: loadsRAM1},
			},
		},
	}
}

type testDeps struct {
	timeGenerator *mock.MockTimeGenerator
	statsRepo     *mock.MockStatsRepo
	store         *mock.MockStore
}

func (deps testDeps) Build(t *testing.T) (*core.Service, error) {
	t.Helper()

	return core.NewService(
		deps.timeGenerator,
		deps.statsRepo,
		deps.store,
	)
}

func initTestDeps(t *testing.T) testDeps {
	t.Helper()

	ctrl := gomock.NewController(t)

	return testDeps{
		timeGenerator: mock.NewMockTimeGenerator(ctrl),
		statsRepo:     mock.NewMockStatsRepo(ctrl),
		store:         mock.NewMockStore(ctrl),
	}
}

const (
	testHookNow                   = "Now"
	testHookGetSensorsByHardware  = "GetSensorsByHardware"
	testHookStoreValue            = "StoreValue"
	testHookGetStatsForRange      = "GetStatsForRange"
	testHookGetCurrentSensorValue = "GetCurrentSensorValue"
)

func initTestHookSet(deps testDeps, data testData) testutils.HookSet {
	hooks := testutils.HookSet{}
	hooks.Add(
		testHookNow,
		deps.timeGenerator.EXPECT().Now(),
		data.now,
	)
	hooks.Add(
		testHookGetSensorsByHardware,
		deps.statsRepo.EXPECT().GetSensorsByHardware(data.ctx),
		data.sensorsByHardware,
		nil,
	)
	hooks.Add(
		testHookGetCurrentSensorValue,
		deps.statsRepo.EXPECT().GetCurrentSensorValues(data.ctx),
		data.currentValuesBySensor,
		nil,
	)
	for sensor, currentSensorValue := range data.currentValuesBySensor {
		value := core.Value{
			Value:     int64(math.Round(float64(currentSensorValue))),
			Timestamp: data.now,
		}

		hooks.Add(
			testHookStoreValue,
			deps.store.EXPECT().StoreValue(sensor.ID, value),
			nil,
		)
	}

	for _, sensorType := range data.valuesPerSensorsForHardware {
		for _, sensors := range sensorType {
			for sensor, values := range sensors {
				hooks.Add(
					testHookGetStatsForRange,
					deps.store.EXPECT().GetValuesForRange(sensor.ID, data.now.Add(-data.req.ForRange), data.now),
					values,
					nil,
				)
			}
		}
	}

	return hooks
}

func generateValuesForRange(from, to time.Time, step time.Duration) []core.Value {
	var values []core.Value

	for t := from; t.Before(to); t = t.Add(step) {
		values = append(values,
			core.Value{
				Value:     rand.N[int64](10) + 40,
				Timestamp: time.UnixMilli(t.UnixMilli()), // to round to milliseconds
			},
		)
	}
	return values
}
