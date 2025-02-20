package storage

import (
	"context"
	"testing"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/internal/utils/rand"
	"github.com/stretchr/testify/require"
)

func setUp(ctx context.Context, t *testing.T) Storage {
	cfg := config.NewWithoutParsing()

	s, err := New(ctx, cfg.DatabaseDSN)
	require.NoError(t, err)

	require.NoError(t, s.InitDB(ctx))

	return s
}

func createRandomUser(ctx context.Context, s Storage, t *testing.T) *User {
	user := NewUser(utils.NewUUID6(), rand.RandomString(10), "somepass")
	err := s.CreateUser(ctx, user)
	require.NoError(t, err)

	return &user
}
