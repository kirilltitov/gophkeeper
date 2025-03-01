package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// WithTransaction opens a transaction, executes a given func with a transaction,
// rolling back the transaction should func return error,
// and returning T result value if func returned no error.
//
// WARNING: one MUST eventually manually commit the transaction in passed func.
func WithTransaction[T any](ctx context.Context, pg *PgSQL, f func(pgx.Tx) (*T, error)) (*T, error) {
	transaction, err := pg.Conn.Begin(ctx)
	if err != nil {
		return nil, err
	}

	result, err := f(transaction)
	if err != nil {
		if err := transaction.Rollback(ctx); err != nil {
			if errors.Is(err, pgx.ErrTxClosed) {
				return result, nil
			} else {
				return nil, err
			}
		}
		return nil, err
	}

	return result, nil
}

// WithVoidTransaction opens a transaction, executes a given func with a transaction,
// rolling back the transaction should func return error,
// and returning void if func returned no error.
//
// WARNING: one MUST eventually manually commit the transaction in passed func.
func WithVoidTransaction(ctx context.Context, pg *PgSQL, f func(pgx.Tx) error) error {
	_, err := WithTransaction(ctx, pg, func(tx pgx.Tx) (*any, error) {
		return nil, f(tx)
	})

	return err
}
