package http

import (
	"fmt"

	openapi "github.com/genvmoroz/win-stats-picker/internal/http/generated"
	"github.com/labstack/echo/v4/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"
)

func (s *Server) register() error {
	// register OpenAPI handlers
	swagger, err := openapi.GetSwagger()
	if err != nil {
		return fmt.Errorf("get swagger: %w", err)
	}
	swagger.Servers = nil
	handler := openapi.NewStrictHandler(s.router, nil)
	openapi.RegisterHandlers(s.echo, handler)
	s.echo.Use(oapimiddleware.OapiRequestValidator(swagger))

	s.echo.Use(middleware.Logger())

	return nil
}
