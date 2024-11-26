package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrCount     = errors.New("count error")
	ErrExecute   = errors.New("execute error")
	ErrSelect    = errors.New("select error")
	ErrSelectOne = errors.New("select (one row) error")
)

type NamedArgs = pgx.NamedArgs

// Count executes the provided SQL expecting a COUNT.
func Count(ctx context.Context, db *Database, sql string, args ...any) (int64, error) {
	db.logger.Debug("query", "sql", sql, "args", args)

	res, err := db.cnx.Query(ctx, sql, args...)
	if err != nil {
		return -1, errors.Join(ErrCount, err)
	}

	defer res.Close()

	count, err := pgx.CollectExactlyOneRow(res, pgx.RowTo[int64])
	if err != nil {
		return -1, errors.Join(ErrCount, err)
	}

	// Rows MUST be closed prior to reading the error.
	// CollectExactlyOneRow does that already.
	if err := res.Err(); err != nil {
		return -1, errors.Join(ErrCount, err)
	}

	return count, nil
}

// Execute executes the provided SQL string without expecting anything to return.
func Execute(ctx context.Context, db *Database, sql string, args ...any) error {
	db.logger.Debug("query", "sql", sql, "args", args)

	res, err := db.cnx.Query(ctx, sql, args...)
	if err != nil {
		return errors.Join(ErrExecute, err)
	}

	defer res.Close()

	if err := res.Err(); err != nil {
		return errors.Join(ErrExecute, err)
	}

	return nil
}

// MustSelectOne executes the provided SQL and return the found row.
// It returns an error if none, or if more than one rows are found.
func MustSelectOne[T any](ctx context.Context, db *Database, sql string, args ...any) (*T, error) {
	db.logger.Debug("query", "sql", sql, "args", args)

	res, err := db.cnx.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Join(ErrSelectOne, err)
	}

	defer res.Close()

	out, err := pgx.CollectExactlyOneRow(res, pgx.RowToStructByNameLax[T])
	if err != nil {
		return nil, errors.Join(ErrSelectOne, err)
	}

	// Rows MUST be closed prior to reading the error.
	// CollectExactlyOneRow does that already.
	if err := res.Err(); err != nil {
		return nil, errors.Join(ErrSelectOne, err)
	}

	return &out, nil
}

// Select executes the provided SQL and returns the whole resultset.
func Select[T any](ctx context.Context, db *Database, sql string, args ...any) ([]T, error) {
	db.logger.Debug("query", "sql", sql, "args", args)

	var out []T

	res, err := db.cnx.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Join(ErrSelect, err)
	}

	defer res.Close()

	out, err = pgx.CollectRows(res, pgx.RowToStructByNameLax[T])
	if err != nil {
		return nil, errors.Join(ErrSelect, err)
	}

	// Rows MUST be closed prior to reading the error.
	// CollectRows does that already.
	if err := res.Err(); err != nil {
		return nil, errors.Join(ErrSelect, err)
	}

	return out, nil
}

// SelectOne executes the provided SQL and return the found row.
// It returns an error if more than one rows are found.
func SelectOne[T any](ctx context.Context, db *Database, sql string, args ...any) (*T, error) {
	res, err := MustSelectOne[T](ctx, db, sql, args...)

	switch {
	case err == nil:
		return res, nil
	case errors.Is(err, pgx.ErrNoRows):
		return nil, nil //nolint:nilnil // As expected.
	default:
		return nil, err
	}
}
