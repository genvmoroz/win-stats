package http

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type (
	Config struct {
		Port uint `envconfig:"APP_HTTP_API_PORT" default:"8080"`
	}

	Server struct {
		cfg    Config
		router *Router
		echo   *echo.Echo
		logger logrus.FieldLogger
	}
)

func NewServer(cfg Config, router *Router, logger logrus.FieldLogger) (*Server, error) {
	if lo.IsNil(router) {
		return nil, fmt.Errorf("router is nil")
	}
	if lo.IsNil(logger) {
		return nil, fmt.Errorf("logger is nil")
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("can't listen on port %d: %w", cfg.Port, err)
	}

	if err = ln.Close(); err != nil {
		return nil, fmt.Errorf("can't close listener: %w", err)
	}

	server := &Server{
		cfg:    cfg,
		router: router,
		echo:   echo.New(),
		logger: logger,
	}

	if err = server.register(); err != nil {
		return nil, fmt.Errorf("register http server: %w", err)
	}

	return server, nil
}

// Run starts the server and listens for incoming requests.
// The server will be stopped when the context is canceled.
func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func(ch chan error) {
		s.logger.Debug("starting http server")
		ch <- s.echo.Start(fmt.Sprintf(":%d", s.cfg.Port))
	}(errChan)

	select {
	case <-ctx.Done():
	case err := <-errChan:
		return err
	}

	timeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.echo.Shutdown(timeout); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}
	s.logger.Debug("http server stopped")

	return nil
}
