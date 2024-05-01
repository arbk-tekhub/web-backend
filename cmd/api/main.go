package main

import (
	"os"
	"runtime/debug"

	"github.com/benk-techworld/www-backend/cmd/api/router"
	"github.com/benk-techworld/www-backend/internal/app"
	"github.com/benk-techworld/www-backend/internal/db"
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

	r := router.Routes()

	db, err := db.Open(true)
	if err != nil {
		return err
	}
	defer db.Close()

	return app.ServeHTTP(r)
}
