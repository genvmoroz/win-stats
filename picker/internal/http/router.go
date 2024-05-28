package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/genvmoroz/win-stats-picker/internal/core"
	openapi "github.com/genvmoroz/win-stats-picker/internal/http/generated"
	"github.com/samber/lo"
)

type Service interface {
	GetStats(ctx context.Context) (core.GetStatsResponse, error)
}

type Router struct {
	srv         Service
	transformer Transformer
}

func NewRouter(srv Service) (*Router, error) {
	if srv == nil {
		return nil, fmt.Errorf("service is nil")
	}
	return &Router{
		srv:         srv,
		transformer: Transformer{},
	}, nil
}

var _ openapi.StrictServerInterface = (*Router)(nil)

func (r *Router) GetStats(
	ctx context.Context,
	_ openapi.GetStatsRequestObject,
) (openapi.GetStatsResponseObject, error) {
	stats, err := r.srv.GetStats(ctx)
	if err != nil {
		return openapi.GetStats500JSONResponse(wrapAPIError(err, http.StatusInternalServerError)), nil
	}

	resp := r.transformer.GetStatsResponseFromCore(stats)

	return openapi.GetStats200JSONResponse(resp), nil
}

func (r *Router) HealthCheck(
	_ context.Context,
	_ openapi.HealthCheckRequestObject,
) (openapi.HealthCheckResponseObject, error) {
	return openapi.HealthCheck200TextResponse("Up and running!"), nil
}

func wrapAPIError(err error, statusCode int) openapi.Error {
	return openapi.Error{
		Message:    lo.ToPtr(err.Error()),
		StatusCode: lo.ToPtr(int32(statusCode)),
	}
}
