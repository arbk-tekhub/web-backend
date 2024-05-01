package app

import (
	"log/slog"
	"os"
	"strings"

	"github.com/benk-techworld/www-backend/internal/env"
	"github.com/lmittmann/tint"
)

var logger *slog.Logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: logLevel()}))

func GetLogger() *slog.Logger {
	return logger
}

func logLevel() slog.Level {

	logLevel := env.GetString("LOG_LEVEL", "debug")
	switch strings.ToLower(logLevel) {
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}

	return slog.LevelDebug
}
