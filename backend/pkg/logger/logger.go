// Package logger provides slog instance for writing logs to stdout
package logger

import (
	"io"
	"log/slog"
	"os"
)

var Log = slog.New(
	slog.NewJSONHandler(
		io.Writer(
			os.Stdout,
		),
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	),
)
