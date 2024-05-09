package autocleanup

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestTaskStart(t *testing.T) {
	defer goleak.VerifyNone(t)

	repo := &testRepo{}

	cfg := Config{
		Interval:       1 * time.Second,
		CleanOlderThan: 1 * time.Second,
	}
	task, err := NewTask(cfg, repo, logrus.New())
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	task.Start(ctx)

	<-ctx.Done()

	require.Equal(t, int8(2), repo.invokedTimes)
}

type testRepo struct {
	invokedTimes int8
}

func (r *testRepo) DeleteOlderValues(_ time.Time) error {
	r.invokedTimes++

	return nil
}
