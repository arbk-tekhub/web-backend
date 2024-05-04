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

// @title     Benk Techworld API
// @version 1.0
// @description An API written in Go using Gin framework.

// @contact.name Arafet BenKilani
// @contact.url https://www.linkedin.com/in/arafet-ben-kilani/
// @contact.email mr.arafetk@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/license/mit

// @host localhost:8080
// @BasePath  /v1

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
	basicAuth struct {
		username       string
		hashedPassword string
	}
}

func run() error {

	var cfg config

	flag.IntVar(&cfg.httpPort, "port", 8080, "http server port")
	flag.Parse()

	cfg.env = env.GetString("APP_ENV", "development")
	cfg.logLevel = env.GetLogLevel("LOG_LEVEL", slog.LevelDebug)
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: cfg.logLevel}))

	cfg.basicAuth.username = env.GetString("BASIC_AUTH_USERNAME", "admin")
	cfg.basicAuth.hashedPassword = env.GetString("BASIC_AUTH_HASHED_PASSWORD", "$2a$12$tzGCcHO3lEWT33elGOhh0uz485PUq.YUiR8U2c98/drCPCsGGtLlu")

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
