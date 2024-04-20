package http

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) setupRoutes() {
	s.logger.Debug("setting up routes")

	s.echo.GET("/stats", s.GetStats)
	s.echo.GET("/health", s.GetHealthcheck)
	s.echo.GET("/metrics", echoprometheus.NewHandler())
	s.echo.Use(echoprometheus.NewMiddleware("http_server"))
	s.echo.Use(middleware.Logger())
	pprof.Register(s.echo)
}
