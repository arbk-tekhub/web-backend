package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/arbk-tekhub/www-backend/cmd/api/router"
	"github.com/arbk-tekhub/www-backend/internal/app"
	"github.com/arbk-tekhub/www-backend/internal/config"
	"github.com/lmittmann/tint"
)

func main() {

	app.SetLogger(&tint.Options{Level: slog.LevelDebug})
	logger := app.GetLogger()

	err := run()
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}

}

func run() error {

	cfgFilePath := flag.String("config", "", "path to configuration file")
	flag.Parse()

	if *cfgFilePath == "" {
		return errors.New("missing config file")
	}

	err := config.Load("json", *cfgFilePath)
	if err != nil {
		return err
	}

	r := router.Routes()

	return app.ServeHTTP(r)
}
