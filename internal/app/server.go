package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/benk-techworld/www-backend/internal/config"
	"github.com/gin-gonic/gin"
)

const (
	defaultAddr           = ":8080"
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

var wg sync.WaitGroup

func ServeHTTP(r *gin.Engine) error {

	cfg := config.Get()

	var addr string

	port := cfg.GetInt("server.port")
	if port == 0 {
		addr = defaultAddr
	} else {
		addr = fmt.Sprintf(":%d", port)
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	}

	shutdownErrorChan := make(chan error)

	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan
		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	logger.Info(fmt.Sprintf("started server on %s", srv.Addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		return err
	}

	logger.Warn("stopped server")

	wg.Wait()

	return nil
}
