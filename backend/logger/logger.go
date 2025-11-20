package logger

import (
	. "backend/config"
	"io"
	"log/slog"
	"os"
)

func initLogger() *slog.Logger {
	logfile, err := os.OpenFile(Cfg.Logs.PathToBackend, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	wrt := io.MultiWriter(os.Stdout, logfile)
	return slog.New(slog.NewJSONHandler(wrt, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

var Logger = initLogger()
