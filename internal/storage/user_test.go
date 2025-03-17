package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/internal/utils/rand"
)

func TestPgSQL_CreateUser(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	var err error
	user1 := createRandomUser(ctx, s, t)

	user2SameLogin := NewUser(utils.NewUUID6(), user1.Login, "someotherpass")
	err = s.CreateUser(ctx, user2SameLogin)
	require.ErrorIs(t, err, ErrDuplicateUserFound)

	user3SameID := NewUser(user1.ID, user1.Login+"2", "someotherotherpass")
	err = s.CreateUser(ctx, user3SameID)
	require.ErrorIs(t, err, ErrDuplicateUserFound)
}

func TestPgSQL_LoadUser(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	var loadedUser *User
	var err error

	loadedUser, err = s.LoadUser(ctx, rand.RandomString(10))
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, loadedUser)

	user := createRandomUser(ctx, s, t)

	loadedUser, err = s.LoadUser(ctx, user.Login)
	require.NoError(t, err)
	// require.True(t, user.CreatedAt.Equal(loadedUser.CreatedAt))
	loadedUser.CreatedAt = user.CreatedAt
	require.Equal(t, user, loadedUser)
}
