package picker

import (
	"context"
	"fmt"

	"github.com/genvmoroz/win-stats-prometheus-collector/internal/core"
	openapi "github.com/genvmoroz/win-stats-prometheus-collector/internal/repository/picker/generated"
	"github.com/hashicorp/go-cleanhttp"
)

type Repo struct {
	client      *openapi.Client
	transformer Transformer
}

func NewRepo(ctx context.Context, host string) (*Repo, error) {
	if host == "" {
		return nil, fmt.Errorf("host is empty")
	}

	opts := []openapi.ClientOption{
		openapi.WithBaseURL(host),
		openapi.WithHTTPClient(cleanhttp.DefaultClient()),
	}

	client, err := openapi.NewClient(host, opts...)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	repo := &Repo{
		client:      client,
		transformer: Transformer{},
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
