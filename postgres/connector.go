package postgres

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luca-arch/go-goodies/logger"
)

const MaxConnectionAttempts = 5

var ErrConnect = errors.New("could not connect")

type Database struct {
	cnx    *pgxpool.Pool
	logger *slog.Logger
}

// Connect instantiates a new connection pool from the provided DSN string.
func Connect(ctx context.Context, dsn string, l *slog.Logger) (*Database, error) {
	var err error

	if l == nil {
		l = logger.NewNop()
	}

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errors.Join(ErrConnect, err)
	}

	for attempt := range MaxConnectionAttempts {
		err = db.Ping(ctx)
		if err == nil {
			return &Database{cnx: db, logger: l}, nil
		}

		if attempt == MaxConnectionAttempts-1 {
			break
		}

		l.Warn("waiting for database", "attempt", attempt+1)

		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return nil, errors.Join(ErrConnect, err)
}

// MustConnect instantiates a new connection pool from the provided DSN string or panics.
func MustConnect(ctx context.Context, dsn string, l *slog.Logger) *Database {
	db, err := Connect(ctx, dsn, l)
	if err != nil {
		panic(err)
	}

	return db
}
