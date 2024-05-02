package main

import (
	"flag"
	"os"
	"runtime/debug"

	"github.com/benk-techworld/www-backend/cmd/api/router"
	"github.com/benk-techworld/www-backend/internal/app"
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

	var port int

	flag.IntVar(&port, "port", 8080, "http server port")
	flag.Parse()

	r := router.Routes()

	return app.ServeHTTP(port, r)
}
