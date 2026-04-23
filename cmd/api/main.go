package main

import (
	"log/slog"
	"os"

	"apigo/internal/config"
	"apigo/internal/pkg/app"
	"apigo/pkg/log/sl"
)

func main() {
	c := config.MustLoad()
	a := app.New(c)

	if err := a.Run(); err != nil {
		slog.Error("application error", sl.Err(err))
		os.Exit(1)
	}
}
