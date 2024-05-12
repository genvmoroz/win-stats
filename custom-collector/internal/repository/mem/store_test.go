package mem

import (
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/araddon/dateparse"
	"github.com/genvmoroz/custom-collector/internal/core"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

const (
	testTemperature core.SensorID = "temperature"
	testFanSpeed    core.SensorID = "fan_speed"
)

func TestStoreStoreAndGetValuesForRange(t *testing.T) {
	t.Parallel()

	type (
		input struct {
			sID  core.SensorID
			from time.Time
			to   time.Time
		}
		want struct {
			resp       []core.Value
			errPresent bool
		}
	)

	tests := []struct {
		name  string
		pre   func(r *Store)
		input input
		want  want
	}{
		{
			name: "store and get values for range",
			pre: func(r *Store) {
				values := []core.Value{
					{Value: 1, Timestamp: dateparse.MustParse("2021-01-01")},
					{Value: 2, Timestamp: dateparse.MustParse("2021-01-02")},
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
					{Value: 5, Timestamp: dateparse.MustParse("2021-01-05")},
				}
				for _, v := range values {
					require.NoError(t, r.StoreValue(testTemperature, v))
				}
			},
			input: input{
				sID:  testTemperature,
				from: dateparse.MustParse("2021-01-02"),
				to:   dateparse.MustParse("2021-01-04"),
			},
			want: want{
				resp: []core.Value{
					{Value: 2, Timestamp: dateparse.MustParse("2021-01-02")},
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
				},
			},
		},
		{
			name: "stored values is not sorted",
			pre: func(r *Store) {
				values := []core.Value{
					{Value: 5, Timestamp: dateparse.MustParse("2021-01-05")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 2, Timestamp: dateparse.MustParse("2021-01-02")},
					{Value: 1, Timestamp: dateparse.MustParse("2021-01-01")},
				}
				for _, v := range values {
					require.NoError(t, r.StoreValue(testTemperature, v))
				}
			},
			input: input{
				sID:  testTemperature,
				from: dateparse.MustParse("2021-01-02"),
				to:   dateparse.MustParse("2021-01-04"),
			},
			want: want{
				resp: []core.Value{
					{Value: 2, Timestamp: dateparse.MustParse("2021-01-02")},
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
				},
			},
		},
		{
			name: "no values",
			input: input{
				sID:  testTemperature,
				from: dateparse.MustParse("2021-01-02"),
				to:   dateparse.MustParse("2021-01-04"),
			},
			want: want{
				resp: nil,
			},
		},
		{
			name: "no records for the range",
			pre: func(r *Store) {
				values := []core.Value{
					{Value: 1, Timestamp: dateparse.MustParse("2021-01-01")},
					{Value: 2, Timestamp: dateparse.MustParse("2021-01-02")},
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
					{Value: 5, Timestamp: dateparse.MustParse("2021-01-05")},
				}
				for _, v := range values {
					require.NoError(t, r.StoreValue(testTemperature, v))
				}
			},
			input: input{
				sID:  testTemperature,
				from: dateparse.MustParse("2021-01-06"),
				to:   dateparse.MustParse("2021-01-07"),
			},
			want: want{
				resp: nil,
			},
		},
		{
			name: "no records for the sensor type",
			pre: func(r *Store) {
				values := []core.Value{
					{Value: 1, Timestamp: dateparse.MustParse("2021-01-01")},
					{Value: 2, Timestamp: dateparse.MustParse("2021-01-02")},
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
					{Value: 5, Timestamp: dateparse.MustParse("2021-01-05")},
				}
				for _, v := range values {
					require.NoError(t, r.StoreValue(testTemperature, v))
				}
			},
			input: input{
				sID:  testFanSpeed,
				from: dateparse.MustParse("2021-01-02"),
				to:   dateparse.MustParse("2021-01-04"),
			},
			want: want{
				resp: nil,
			},
		},
		{
			name: "from time is after the to time",
			input: input{
				sID:  testTemperature,
				from: dateparse.MustParse("2021-01-02"),
				to:   dateparse.MustParse("2021-01-01"),
			},
			want: want{
				resp:       nil,
				errPresent: true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := NewStore(logrus.New())
			require.NoError(t, err)
			require.NotNil(t, r)

			if tt.pre != nil {
				tt.pre(r)
			}

			got, err := r.GetValuesForRange(
				tt.input.sID,
				tt.input.from,
				tt.input.to,
			)
			require.Equal(t, tt.want.errPresent, err != nil)
			compareValues(t, tt.want.resp, got)
		})
	}
}

func TestStoreDeleteOlderValues(t *testing.T) {
	t.Parallel()

	type (
		input struct {
			t time.Time
		}
		want struct {
			errPresent bool
		}
	)
	tests := []struct {
		name   string
		pre    func(r *Store)
		input  input
		want   want
		assert func(t *testing.T, r *Store)
	}{
		{
			name: "success",
			pre: func(r *Store) {
				values := []core.Value{
					{Value: 1, Timestamp: dateparse.MustParse("2021-01-01")},
					{Value: 2, Timestamp: dateparse.MustParse("2021-01-02")},
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
					{Value: 5, Timestamp: dateparse.MustParse("2021-01-05")},
				}
				for _, v := range values {
					require.NoError(t, r.StoreValue(testTemperature, v))
				}
			},
			input: input{
				t: dateparse.MustParse("2021-01-03"),
			},
			want: want{
				errPresent: false,
			},
			assert: func(t *testing.T, r *Store) {
				got, err := r.GetValuesForRange(testTemperature, time.Time{}, dateparse.MustParse("2021-01-06"))
				require.NoError(t, err)

				exp := []core.Value{
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
					{Value: 5, Timestamp: dateparse.MustParse("2021-01-05")},
				}
				compareValues(t, exp, got)
			},
		},
		{
			name: "no values",
			input: input{
				t: dateparse.MustParse("2021-01-03"),
			},
			want: want{
				errPresent: false,
			},
			assert: func(t *testing.T, r *Store) {
				got, err := r.GetValuesForRange(
					testTemperature,
					time.Time{},
					dateparse.MustParse("2021-01-06"),
				)
				require.NoError(t, err)
				require.Empty(t, got)
			},
		},
		{
			name: "all deleted",
			pre: func(r *Store) {
				values := []core.Value{
					{Value: 1, Timestamp: dateparse.MustParse("2021-01-01")},
					{Value: 2, Timestamp: dateparse.MustParse("2021-01-02")},
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
					{Value: 5, Timestamp: dateparse.MustParse("2021-01-05")},
				}
				for _, v := range values {
					require.NoError(t, r.StoreValue(testTemperature, v))
				}
			},
			input: input{
				t: dateparse.MustParse("2021-01-06"),
			},
			want: want{
				errPresent: false,
			},
			assert: func(t *testing.T, r *Store) {
				got, err := r.GetValuesForRange(testTemperature, time.Time{}, dateparse.MustParse("2021-01-10"))
				require.NoError(t, err)
				require.Empty(t, got)
			},
		},
		{
			name: "no older values",
			pre: func(r *Store) {
				values := []core.Value{
					{Value: 1, Timestamp: dateparse.MustParse("2021-01-01")},
					{Value: 2, Timestamp: dateparse.MustParse("2021-01-02")},
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
					{Value: 5, Timestamp: dateparse.MustParse("2021-01-05")},
				}
				for _, v := range values {
					require.NoError(t, r.StoreValue(testTemperature, v))
				}
			},
			input: input{
				t: dateparse.MustParse("2021-01-01"),
			},
			want: want{
				errPresent: false,
			},
			assert: func(t *testing.T, r *Store) {
				got, err := r.GetValuesForRange(testTemperature, time.Time{}, dateparse.MustParse("2021-01-10"))
				require.NoError(t, err)
				exp := []core.Value{
					{Value: 1, Timestamp: dateparse.MustParse("2021-01-01")},
					{Value: 2, Timestamp: dateparse.MustParse("2021-01-02")},
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
					{Value: 5, Timestamp: dateparse.MustParse("2021-01-05")},
				}
				compareValues(t, exp, got)
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			r, err := NewStore(logrus.New())
			require.NoError(t, err)
			require.NotNil(t, r)

			if tt.pre != nil {
				tt.pre(r)
			}

			err = r.DeleteOlderValues(tt.input.t)
			require.Equal(t, tt.want.errPresent, err != nil)

			if tt.assert != nil {
				tt.assert(t, r)
			}
		})
	}
}

func TestStoreClose(t *testing.T) {
	t.Parallel()

	type want struct {
		errPresent bool
	}
	tests := []struct {
		name   string
		pre    func(r *Store)
		want   want
		assert func(t *testing.T, r *Store)
	}{
		{
			name: "success",
			pre: func(r *Store) {
				values := []core.Value{
					{Value: 1, Timestamp: dateparse.MustParse("2021-01-01")},
					{Value: 2, Timestamp: dateparse.MustParse("2021-01-02")},
					{Value: 3, Timestamp: dateparse.MustParse("2021-01-03")},
					{Value: 4, Timestamp: dateparse.MustParse("2021-01-04")},
					{Value: 5, Timestamp: dateparse.MustParse("2021-01-05")},
				}
				for _, v := range values {
					require.NoError(t, r.StoreValue(testTemperature, v))
				}
			},
			want: want{
				errPresent: false,
			},
			assert: func(t *testing.T, r *Store) {
				require.Empty(t, r.dbs)
				got, err := r.GetValuesForRange(testTemperature, time.Time{}, dateparse.MustParse("2021-01-10"))
				require.NoError(t, err)
				require.Empty(t, got)
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r, err := NewStore(logrus.New())
			require.NoError(t, err)

			if tt.pre != nil {
				tt.pre(r)
			}

			err = r.Close()
			require.NoError(t, err)

			<-time.After(10 * time.Millisecond) // wait for the goroutine to close the dbs

			if tt.assert != nil {
				tt.assert(t, r)
			}
		})
	}
}

func TestStoreFlowInParallel(t *testing.T) {
	defer goleak.VerifyNone(t)

	from := time.Now()
	to := time.Now().Add(1 * time.Hour)

	// Generate values
	var (
		cpuTemps = generateValuesForRange(from, to, time.Second)
		cpuFans  = generateValuesForRange(from, to, time.Second)
	)

	// Create Store
	store, err := NewStore(logrus.New())
	require.NoError(t, err)

	// Store values in parallel
	wg := &sync.WaitGroup{}

	type actionFunc func(t *testing.T, wg *sync.WaitGroup, value core.Value)

	foreachInParallel := func(t *testing.T, wg *sync.WaitGroup, values []core.Value, action actionFunc) {
		defer wg.Done()
		for _, v := range values {
			wg.Add(1)
			go action(t, wg, v)
		}
	}

	storeFunc := func(sID core.SensorID) actionFunc {
		return func(t *testing.T, wg *sync.WaitGroup, value core.Value) {
			defer wg.Done()
			require.NoError(t, store.StoreValue(sID, value))
		}
	}

	wg.Add(1)
	go foreachInParallel(t, wg, cpuTemps, storeFunc(testTemperature))
	wg.Add(1)
	go foreachInParallel(t, wg, cpuFans, storeFunc(testFanSpeed))

	wg.Wait()

	// Get values for range in parallel
	getAndCompare := func(
		t *testing.T,
		wg *sync.WaitGroup,
		sID core.SensorID,
		expected []core.Value,
	) {
		defer wg.Done()

		got, err := store.GetValuesForRange(sID, from, to)
		require.NoError(t, err)
		compareValues(t, expected, got)
	}
	wg.Add(1)
	go getAndCompare(t, wg, testTemperature, cpuTemps)

	wg.Add(1)
	go getAndCompare(t, wg, testFanSpeed, cpuFans)

	wg.Wait()

	err = store.Close()
	require.NoError(t, err)
}

func compareValues(t *testing.T, exp, got []core.Value) {
	t.Helper()

	require.Len(t, got, len(exp))
	for i, exp := range exp {
		require.Equal(t, exp.Value, got[i].Value)
		require.True(t, exp.Timestamp.Equal(got[i].Timestamp))
	}
}

func generateValuesForRange(from, to time.Time, step time.Duration) []core.Value {
	var values []core.Value

	for t := from; t.Before(to); t = t.Add(step) {
		values = append(values,
			core.Value{
				Value:     rand.N[int64](100),
				Timestamp: time.UnixMilli(t.UnixMilli()), // to round to milliseconds
			},
		)
	}
	return values
}
