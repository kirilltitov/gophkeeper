package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/internal/utils/rand"
)

func setUp(ctx context.Context, t *testing.T) Storage {
	cfg := config.NewWithoutParsing()

	s, err := New(ctx, cfg.DatabaseDSN)
	require.NoError(t, err)

	err = s.InitDB(ctx)
	var connError *pgconn.ConnectError
	if errors.As(err, &connError) {
		t.Skipf("Could not connect to PgSQL, skipping all storage tests")
		return nil
	}

	require.NoError(t, err)

	return s
}

func createRandomUser(ctx context.Context, s Storage, t *testing.T) *User {
	user := NewUser(utils.NewUUID6(), rand.RandomString(10), "somepass")
	err := s.CreateUser(ctx, user)
	require.NoError(t, err)

	return &user
}
