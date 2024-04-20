package autocleanup

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type (
	Store interface {
		DeleteOlderValues(t time.Time) error
	}

	Config struct {
		Interval       time.Duration `envconfig:"APP_AUTO_CLEANUP_INTERVAL" default:"1h"`
		CleanOlderThan time.Duration `envconfig:"APP_AUTO_CLEANUP_OLDER_THAN" default:"1h"`
	}

	Task struct {
		store     Store
		scheduler *gocron.Scheduler

		logger logrus.FieldLogger
	}
)

func NewTask(cfg Config, store Store, logger logrus.FieldLogger) (*Task, error) {
	if lo.IsNil(store) {
		return nil, errors.New("store is nil")
	}
	if lo.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}
	if cfg.Interval <= 0 {
		return nil, errors.New("interval must be greater than 0")
	}
	if cfg.CleanOlderThan <= 0 {
		return nil, errors.New("clean older than param must be greater than 0")
	}

	task := &Task{
		store:     store,
		scheduler: gocron.NewScheduler(),
		logger:    logger,
	}

	err := task.scheduler.
		Every(uint64(math.Round(cfg.Interval.Seconds()))).
		Seconds().
		Do(task.clean(cfg.CleanOlderThan))
	if err != nil {
		return nil, fmt.Errorf("schedule cleanup task: %w", err)
	}

	return task, nil
}

// Start starts the task and runs it until the context is canceled.
func (task *Task) Start(ctx context.Context) {
	stopChan := task.scheduler.Start()
	task.logger.Info("auto cleanup task started")

	<-ctx.Done()

	stopChan <- true
	task.scheduler.Clear()
	task.logger.Info("auto cleanup task stopped")
}

func (task *Task) clean(olderThan time.Duration) func() {
	return func() {
		task.logger.Infof("cleaning stats older than %s\n", olderThan)
		if err := task.store.DeleteOlderValues(time.Now().Add(-1 * olderThan)); err != nil {
			task.logger.Errorf("failed to clean old stats: %v", err)
		}
	}
}
