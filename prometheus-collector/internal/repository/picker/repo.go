package picker

import (
	"context"
	"fmt"
	"time"

	"github.com/genvmoroz/win-stats-prometheus-collector/internal/core"
	openapi "github.com/genvmoroz/win-stats-prometheus-collector/internal/repository/picker/generated"
	"github.com/hashicorp/go-cleanhttp"
)

type (
	Config struct {
		Host    string        `envconfig:"APP_PICKER_REPO_HOST" validate:"required"`
		Timeout time.Duration `envconfig:"APP_PICKER_REPO_TIMEOUT" default:"5m"`
	}

	Repo struct {
		client      *openapi.Client
		transformer Transformer
	}
)

func NewRepo(ctx context.Context, cfg Config) (*Repo, error) {
	if cfg.Host == "" {
		return nil, fmt.Errorf("host is empty")
	}
	if cfg.Timeout <= 0 {
		return nil, fmt.Errorf("timeout must be greater than 0")
	}

	baseClient := cleanhttp.DefaultClient()
	baseClient.Timeout = cfg.Timeout

	opts := []openapi.ClientOption{
		openapi.WithBaseURL(cfg.Host),
		openapi.WithHTTPClient(baseClient),
	}

	client, err := openapi.NewClient(cfg.Host, opts...)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	repo := &Repo{
		client:      client,
		transformer: Transformer{},
	}

	if err = repo.HealthCheck(ctx); err != nil {
		return nil, fmt.Errorf("health check: %w", err)
	}

	return repo, nil
}

func (r *Repo) GetStats(ctx context.Context) (core.Stats, error) {
	resp, err := r.client.GetStats(ctx)
	if err != nil {
		return core.Stats{}, fmt.Errorf("get stats: %w", err)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	stats, err := handleResponse(resp, openapi.ParseGetStatsResponse, r.transformer.GetStatsResponseFromOpenAPI)
	if err != nil {
		return core.Stats{}, err
	}

	return stats, nil
}

func (r *Repo) HealthCheck(ctx context.Context) error {
	resp, err := r.client.HealthCheck(ctx)
	if err != nil {
		return fmt.Errorf("health check: %w", err)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	_, err = handleResponse(resp, openapi.ParseHealthCheckResponse, r.transformer.HealthCheckFromOpenAPI)
	if err != nil {
		return err
	}

	return nil
}
