package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/benk-techworld/www-backend/cmd/api/router"
	"github.com/benk-techworld/www-backend/internal/app"
	"github.com/benk-techworld/www-backend/internal/db"
	"github.com/benk-techworld/www-backend/internal/env"
)

func main() {

	logger := app.GetLogger()

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}

}

func run(logger *slog.Logger) error {

	var port int

	flag.IntVar(&port, "port", 8080, "http server port")
	flag.Parse()

	uri := env.GetString("MONGO_URI", "mongodb://localhost:27017")
	client, err := db.Open(uri)
	if err != nil {
		return err
	}

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	logger.Info("connection to the database has been established successfully")

	r := router.Routes()

	return app.ServeHTTP(port, r)
}
