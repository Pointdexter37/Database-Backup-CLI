package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func InitLogger(logFile string) error {
	var handler slog.Handler

	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		handler = slog.NewJSONHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)

	return nil
}
