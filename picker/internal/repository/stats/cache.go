package stats

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/genvmoroz/win-stats/picker/internal/core"
	"github.com/samber/lo"
)

type CachedRepoConfig struct {
	Retention time.Duration `envconfig:"APP_STATS_CACHE_RETENTION" default:"1s"`
}

type CachedRepo struct {
	baseRepo core.StatsRepo
	cache    cache
	mux      *sync.Mutex
}

func NewCachedRepo(baseRepo core.StatsRepo, cfg CachedRepoConfig) (*CachedRepo, error) {
	if lo.IsNil(baseRepo) {
		return nil, errors.New("base repo is nil")
	}

	return &CachedRepo{
		baseRepo: baseRepo,
		cache:    newCache(cfg.Retention),
		mux:      &sync.Mutex{},
	}, nil
}

// todo: implement a test for this method, try to use RWMutex instead of Mutex
func (c *CachedRepo) GetSensorsByHardware(ctx context.Context) (map[core.Hardware][]core.Sensor, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.cache.isValid() {
		return c.cache.get(), nil
	}

	sensorsByHardware, err := c.baseRepo.GetSensorsByHardware(ctx)
	if err != nil {
		return nil, err
	}

	c.cache.set(sensorsByHardware)

	return c.cache.get(), nil
}

type cache struct {
	stats     map[core.Hardware][]core.Sensor
	retention time.Duration
	updatedAt time.Time
}

func newCache(retention time.Duration) cache {
	return cache{
		stats:     make(map[core.Hardware][]core.Sensor),
		retention: retention,
		updatedAt: time.Time{},
	}
}

func (c *cache) isValid() bool {
	return time.Since(c.updatedAt) < c.retention
}

func (c *cache) set(sensorsByHardware map[core.Hardware][]core.Sensor) {
	c.stats = sensorsByHardware
	c.updatedAt = time.Now()
}

func (c *cache) get() map[core.Hardware][]core.Sensor {
	return c.stats
}
