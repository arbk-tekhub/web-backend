package main

import (
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/arbk-tekhub/www-backend/cmd/api/router"
)

func main() {

	r := router.Routes()

	err := r.Run(":8080")
	if err != nil {
		trace := string(debug.Stack())
		slog.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}

}
