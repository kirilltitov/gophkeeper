package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func WithTransaction[T any](ctx context.Context, pg PgSQL, f func(pgx.Tx) (*T, error)) (*T, error) {
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

func WithVoidTransaction(ctx context.Context, pg PgSQL, f func(pgx.Tx) error) error {
	_, err := WithTransaction(ctx, pg, func(tx pgx.Tx) (*any, error) {
		return nil, f(tx)
	})

	return err
}
