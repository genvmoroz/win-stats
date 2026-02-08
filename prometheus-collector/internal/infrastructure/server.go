// Package infrastructure contains infrastructure-layer components like HTTP server.
package infrastructure

import (
	"context"
	"fmt"
	"net"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
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

	server.echo.Use(middleware.RequestLogger())

	return server, nil
}

// Run starts the server and listens for incoming requests.
// The server will be stopped when the context is canceled.
func (s *Server) Run(ctx context.Context) error {
	startCfg := echo.StartConfig{Address: fmt.Sprintf(":%d", s.cfg.ServerPort)}
	return startCfg.Start(ctx, s.echo)
}
