package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"apigo/internal/config"
	"apigo/internal/service/logset"
	chstorage "apigo/internal/storage/clickhouse"
	"apigo/internal/version"
	"apigo/pkg/log"

	httpserver "apigo/internal/transport/http"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type App struct {
	config *config.Config
}

func New(config *config.Config) *App {
	return &App{config}
}

func (a *App) Run() error {
	log.InitDefault(a.config.Env)

	slog.Info("server starting...", slog.Group("revision", version.CommitAttr, version.BranchAttr))

	ch, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{a.config.ClickHouse.Address},
		Auth: clickhouse.Auth{
			Database: a.config.ClickHouse.Database,
			Username: a.config.ClickHouse.Username,
			Password: a.config.ClickHouse.Password,
		},
	})
	if err != nil {
		return fmt.Errorf("clickhouse open: %w", err)
	}
	defer ch.Close() //nolint:errcheck

	if err := ch.Ping(context.Background()); err != nil {
		return fmt.Errorf("clickhouse ping: %w", err)
	}

	slog.Info("connected to clickhouse", slog.String("address", a.config.ClickHouse.Address))

	logsetStorage := chstorage.NewLogsetStorage(ch)
	logsetService := logset.NewService(logsetStorage)
	server := httpserver.NewServer(a.config, logsetService)

	errch := make(chan error, 1)
	go func() {
		errch <- server.Run()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errch:
		return fmt.Errorf("server failed: %w", err)
	case <-quit:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	return nil
}
