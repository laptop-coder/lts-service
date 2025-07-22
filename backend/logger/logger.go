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
	wrt := io.MultiWriter(os.Stdout, logfile)
	return slog.New(slog.NewJSONHandler(wrt, nil))
}

var Logger = initLogger()
