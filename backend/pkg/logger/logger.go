// Package logger provides aliases for writing logs to stdout
package logger

import (
	"io"
	"log/slog"
	"os"
)

// slog instance
var log = slog.New(
	slog.NewJSONHandler(
		io.Writer(
			os.Stdout,
		),
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	),
)


// Aliases

// Debug is alias for writing logs to stdout at the debug (-4) level
func Debug (message string) {
	log.Debug(message)
}

// Info is alias for writing logs to stdout at the info (0) level
func Info (message string) {
	log.Info(message)
}

// Warn is alias for writing logs to stdout at the warning (4) level
func Warn (message string) {
	log.Warn(message)
}

// Error is alias for writing logs to stdout at the error (8) level
func Error (message string) {
	log.Error(message)
}
