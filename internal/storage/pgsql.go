package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kirilltitov/gophkeeper/internal/utils"
)

type PgSQL struct {
	Conn *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*PgSQL, error) {
	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	conf.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}

	utils.Log.Infof("Connected to PgSQL with DSN %s", dsn)

	return &PgSQL{Conn: pool}, nil
}

func (s PgSQL) Close() {
	s.Conn.Close()
}
