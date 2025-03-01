package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// PgSQL is an implementation of Storage interface which loads actual records in PostgreSQL DB.
type PgSQL struct {
	Conn *pgxpool.Pool // Conn is a PostgreSQL connection pool.
}

// New creates and returns a fully configured PgSQL instance.
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

// Close closes all pool connections.
func (s *PgSQL) Close() {
	utils.Log.Infof("Closing PgSQL connection")
	s.Conn.Close()
	utils.Log.Infof("Closed PgSQL connection")
}
