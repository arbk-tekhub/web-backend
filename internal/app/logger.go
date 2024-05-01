package app

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

var logger *slog.Logger

func SetLogger(options *tint.Options) {

	logger = slog.New(tint.NewHandler(os.Stdout, options))
}

func GetLogger() *slog.Logger {
	return logger
}
