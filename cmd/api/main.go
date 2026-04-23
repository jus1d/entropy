package main

import (
	"log/slog"
	"os"

	"entropy/internal/config"
	"entropy/internal/pkg/app"
	"entropy/pkg/log/sl"
)

func main() {
	c := config.MustLoad()
	a := app.New(c)

	if err := a.Run(); err != nil {
		slog.Error("application error", sl.Err(err))
		os.Exit(1)
	}
}
