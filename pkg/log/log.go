package log

import (
	"fmt"
	"log/slog"
	"os"

	"apigo/internal/config"
	"apigo/pkg/log/prettyslog"
)

func InitDefault(env config.Env) {
	var logger *slog.Logger

	switch env {
	case config.EnvLocal:
		logger = prettyslog.Init()
	case config.EnvDevelopment:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.LevelDebug,
			AddSource:   true,
			ReplaceAttr: replaceAttr,
		}))
	case config.EnvProduction:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			AddSource:   true,
			ReplaceAttr: replaceAttr,
		}))
	}

	logger = logger.With(
		slog.String("env", string(env)),
	)

	slog.SetDefault(logger)
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		a.Key = "timestamp"
	}

	if a.Key == slog.SourceKey {
		src, ok := a.Value.Any().(*slog.Source)
		if ok {
			a.Value = slog.GroupValue(
				slog.String("loc", fmt.Sprintf("%s:%d", src.File, src.Line)),
				slog.String("fn", src.Function),
			)
		}
	}

	return a
}
