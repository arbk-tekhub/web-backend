package main

import (
	"flag"
	"os"
	"runtime/debug"

	"github.com/benk-techworld/www-backend/cmd/api/router"
	"github.com/benk-techworld/www-backend/internal/app"
	"github.com/benk-techworld/www-backend/internal/db"
	"github.com/benk-techworld/www-backend/internal/env"
)

func main() {

	logger := app.GetLogger()

	err := run()
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}

}

func run() error {

	var httpPort int
	flag.IntVar(&httpPort, "port", 8080, "http server port")
	flag.Parse()

	dsn := env.GetString("DB_DSN", "")
	automigrate := env.GetBool("DB_AUTOMIGRATE", true)

	db, err := db.Open(dsn, automigrate)
	if err != nil {
		return err
	}
	defer db.Close()

	r := router.Routes()

	return app.ServeHTTP(httpPort, r)
}
