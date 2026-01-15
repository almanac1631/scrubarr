package scrubarr

import (
	"log/slog"
	"os"
)

func setupLogging() {
	var level slog.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		panic(err)
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
}
