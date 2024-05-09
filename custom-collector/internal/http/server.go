package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/genvmoroz/custom-collector/internal/core"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

type (
	Service interface {
		GetStats(ctx context.Context, req core.GetStatsRequest) (core.GetStatsResponse, error)
	}

	Config struct {
		Port                uint          `envconfig:"APP_HTTP_API_PORT" default:"8080"`
		RequestTimeout      time.Duration `envconfig:"APP_HTTP_API_REQUEST_TIMEOUT" default:"10s"`
		StatsUpdateInterval time.Duration `envconfig:"APP_HTTP_API_STATS_UPDATE_INTERVAL" default:"1s"`
	}

	Server struct {
		srv    Service
		echo   *echo.Echo
		logger logrus.FieldLogger

		port                uint
		requestTimeout      time.Duration
		statsUpdateInterval time.Duration
		wsu                 websocket.Upgrader
	}
)

func NewServer(cfg Config, srv Service, logger logrus.FieldLogger) (*Server, error) {
	if lo.IsNil(srv) {
		return nil, fmt.Errorf("service is nil")
	}
	if lo.IsNil(logger) {
		return nil, fmt.Errorf("logger is nil")
	}
	if cfg.RequestTimeout <= 0 {
		return nil, fmt.Errorf("request duration must be more than 0")
	}
	if cfg.StatsUpdateInterval <= 0 {
		return nil, fmt.Errorf("stats update interval must be more than 0")
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("can't listen on port %d: %w", cfg.Port, err)
	}

	if err = ln.Close(); err != nil {
		return nil, fmt.Errorf("can't close listener: %w", err)
	}

	server := &Server{
		srv:                 srv,
		echo:                echo.New(),
		logger:              logger,
		port:                cfg.Port,
		requestTimeout:      cfg.RequestTimeout,
		statsUpdateInterval: cfg.StatsUpdateInterval,
		wsu:                 websocket.Upgrader{},
	}

	server.setupRoutes()

	return server, nil
}

// Run starts the server and listens for incoming requests.
// The server will be stopped when the context is canceled.
func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func(ch chan error) {
		s.logger.Debug("starting http server")
		ch <- s.echo.Start(fmt.Sprintf(":%d", s.port))
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

func (s *Server) GetStats(c echo.Context) error {
	return HandleWS[GetStatsRequest, core.GetStatsRequest](
		c, s.wsu, s.logger,
		parseGetStatsRequest,
		toCoreGetStatsRequest,
		s.getStats,
	)
}

func (s *Server) getStats(ctx context.Context, conn *websocket.Conn, req core.GetStatsRequest) error {
	ticker := time.NewTicker(s.statsUpdateInterval)
	for {
		select {
		case <-ctx.Done():
			s.writeMessageWS(conn, websocket.CloseMessage, "close ws connection")
			break
		default:
		}

		resp, err := s.processGetStatsWithContext(ctx, req)
		if err != nil {
			s.writeMessageWS(conn, websocket.TextMessage, fmt.Sprintf("internal error: %s", err.Error()))
			break
		}
		if err = conn.WriteJSON(resp); err != nil {
			s.writeMessageWS(conn, websocket.TextMessage, fmt.Sprintf("write response: %s", err.Error()))
			break
		}

		select {
		case <-ctx.Done():
		case <-ticker.C:
		}
	}

	return nil
}

func HandleWS[APIReq any, CoreReq any](
	c echo.Context,
	wsu websocket.Upgrader,
	logger logrus.FieldLogger,
	parseReq func(c echo.Context) (APIReq, error),
	toCoreReq func(req APIReq) (CoreReq, error),
	handle func(ctx context.Context, conn *websocket.Conn, req CoreReq) error,
) error {
	if !c.IsWebSocket() {
		return c.String(http.StatusBadRequest, "the endpoint accepts only websocket connections")
	}

	req, err := parseReq(c)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("parse request: %s", err.Error()))
	}

	coreReq, err := toCoreReq(req)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("convert request: %s", err.Error()))
	}

	conn, err := wsu.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("upgrade to websocket: %w", err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			logger.Errorf("close websocket connection error: %s", closeErr.Error())
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// set close handler to cancel the context when the connection is closed by the peer
	setDefaultCloseHandler(conn, cancel, logger)

	// run the default connection reader in a separate goroutine.
	// it will read messages from the peer and log them.
	// It will stop when the context is canceled or the connection is closed.
	go runDefaultConnectionReader(ctx, conn, logger)

	if err = handle(ctx, conn, coreReq); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("handle request: %s", err.Error()))
	}

	return nil
}

func (s *Server) processGetStatsWithContext(ctx context.Context, req core.GetStatsRequest) (GetStatsResponse, error) {
	var zero GetStatsResponse

	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	resp, err := s.srv.GetStats(ctx, req)
	if err != nil {
		return zero, fmt.Errorf("get stats: %w", err)
	}

	return fromCoreResp(resp), nil
}

func (s *Server) GetHealthcheck(c echo.Context) error {
	return c.String(http.StatusOK, "Up and running!")
}
