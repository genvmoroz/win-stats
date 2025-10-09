// Package infrastructure contains infrastructure-layer components like HTTP server.
package infrastructure

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type Config struct {
	ServerPort uint `envconfig:"APP_INFRA_SERVER_PORT" default:"8888"`
}

type Server struct {
	cfg    Config
	echo   *echo.Echo
	logger *zap.SugaredLogger
}

func NewServer(cfg Config, logger *zap.SugaredLogger) (*Server, error) {
	if cfg.ServerPort == 0 {
		return nil, fmt.Errorf("server port is not set")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.ServerPort))
	if err != nil {
		return nil, fmt.Errorf("can't listen on port %d: %w", cfg.ServerPort, err)
	}

	if err = ln.Close(); err != nil {
		return nil, fmt.Errorf("can't close listener: %w", err)
	}

	server := &Server{
		echo:   echo.New(),
		logger: logger,
		cfg:    cfg,
	}

	// register prometheus and pprof handlers
	server.echo.GET("/metrics", echoprometheus.NewHandler())
	pprof.Register(server.echo)

	server.echo.Use(middleware.Logger())

	return server, nil
}

// Run starts the server and listens for incoming requests.
// The server will be stopped when the context is canceled.
func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func(ch chan error) {
		s.logger.Debug("starting infra http server")
		ch <- s.echo.Start(fmt.Sprintf(":%d", s.cfg.ServerPort))
	}(errChan)

	select {
	case <-ctx.Done():
	case err := <-errChan:
		return err
	}

	const shutdownTimeout = 2 * time.Second

	timeout, cancel := context.WithTimeout(context.Background(), shutdownTimeout) //nolint:contextcheck // false-positive: https://github.com/kkHAIKE/contextcheck/issues/2
	defer cancel()

	if err := s.echo.Shutdown(timeout); err != nil {
		return fmt.Errorf("shutdown infra http server: %w", err)
	}
	s.logger.Debug("infra http server stopped")

	return nil
}
