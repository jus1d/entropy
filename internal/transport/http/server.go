package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"apigo/internal/config"
	"apigo/internal/transport/http/middleware"
	v1 "apigo/internal/transport/http/v1"
	"apigo/pkg/ratelimit"
	"apigo/pkg/requestid"
	"apigo/pkg/requestlog"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	http *http.Server
}

func NewServer(c *config.Config) *Server {
	e := echo.New()

	e.HTTPErrorHandler = HTTPErrorHandler

	e.Use(middleware.Recover)
	e.Use(requestid.New)
	e.Use(middleware.BodyLimit("4M"))
	e.Use(ratelimit.TokenBucket(10, 20))
	e.Use(middleware.Timeout(c.Server.Timeout))
	e.Use(middleware.Metrics)
	e.Use(requestlog.Completed)
	e.Pre(echomw.RemoveTrailingSlash())

	switch c.Env {
	case config.EnvLocal, config.EnvDevelopment:
		e.Use(middleware.CORS)
	}

	e.GET("/liveness", liveness)
	e.GET("/readiness", readiness)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	api := e.Group("/api")
	v1Group := api.Group("/v1")
	v1Router := v1.NewRouter()
	v1Router.Register(v1Group)

	return &Server{
		http: &http.Server{
			Addr:         c.Server.Address,
			Handler:      e,
			ReadTimeout:  c.Server.Timeout,
			WriteTimeout: c.Server.Timeout,
			IdleTimeout:  c.Server.IdleTimeout,
		},
	}
}

// Handler returns the underlying http.Handler, useful for testing.
func (s *Server) Handler() http.Handler {
	return s.http.Handler
}

func (s *Server) Run() error {
	slog.Info("server started", slog.String("address", s.http.Addr))

	if err := s.http.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("shutting down...")

	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	slog.Info("server stopped")
	return nil
}
