package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"github.com/benk-techworld/www-backend/internal/data"
	"github.com/benk-techworld/www-backend/internal/db"
	"github.com/benk-techworld/www-backend/internal/env"
	"github.com/benk-techworld/www-backend/internal/service"
	"github.com/lmittmann/tint"
)

func main() {

	err := run()
	if err != nil {
		trace := string(debug.Stack())
		slog.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}

}

type application struct {
	config  config
	logger  *slog.Logger
	service *service.Service
	wg      sync.WaitGroup
}

type config struct {
	httpPort int
	env      string
	logLevel slog.Level
	database struct {
		uri  string
		name string
	}
}

func run() error {

	var cfg config

	flag.IntVar(&cfg.httpPort, "port", 8080, "http server port")
	flag.Parse()

	cfg.env = env.GetString("APP_ENV", "development")
	cfg.logLevel = env.GetLogLevel("LOG_LEVEL", slog.LevelDebug)
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: cfg.logLevel}))

	cfg.database.uri = env.GetString("MONGO_URI", "mongodb://localhost:27017")
	cfg.database.name = env.GetString("MONGO_DB_NAME", "test")

	db, err := db.Open(cfg.database.uri, cfg.database.name)
	if err != nil {
		return err
	}

	defer func() {
		if err = db.Client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	logger.Info("connection to the database has been established successfully")

	repo := data.NewRepo(db)

	app := &application{
		config:  cfg,
		logger:  logger,
		service: service.New(repo),
	}

	return app.serveHTTP()
}
