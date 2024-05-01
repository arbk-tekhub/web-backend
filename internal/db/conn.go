package db

import (
	"context"
	"errors"
	"time"

	"github.com/benk-techworld/www-backend/assets"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func Open(dsn string, automigrate bool) (*DB, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgxPool, err := pgxpool.New(ctx, "postgres://"+dsn)
	if err != nil {
		return nil, err
	}

	if err = pgxPool.Ping(ctx); err != nil {
		return nil, err
	}

	if automigrate {
		iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
		if err != nil {
			return nil, err
		}
		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, "postgres://"+dsn)
		if err != nil {
			return nil, err
		}

		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		default:
			return nil, err
		}
	}

	return &DB{pgxPool}, nil

}
