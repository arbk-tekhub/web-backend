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

	"github.com/benk-techworld/www-backend/internal/env"
	"github.com/gin-gonic/gin"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

var wg sync.WaitGroup

func ServeHTTP(r *gin.Engine) error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.GetInt("APP_PORT", 8080)),
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
