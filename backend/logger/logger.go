package logger

import (
	"io"
	"log/slog"
	"os"
)

func initLogger() *slog.Logger {
	const PATH_TO_LOG string = "./backend.log" // TODO: move const to .env
	logfile, err := os.OpenFile(PATH_TO_LOG, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	logLevel := new(slog.LevelVar)
	if os.Getenv("LTS_SERVICE_DEV_MODE") == "true" {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}

	wrt := io.MultiWriter(os.Stdout, logfile)
	return slog.New(slog.NewJSONHandler(wrt, &slog.HandlerOptions{
		Level: logLevel,
	}))
}

var Logger = initLogger()
