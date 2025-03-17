package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/internal/utils/rand"
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

func TestPgSQL_AddTag(t *testing.T) {
	var err error
	var loadedSecret *Secret

	ctx := context.Background()
	s := setUp(ctx, t)

	secret := createRandomSecret(t, ctx, s)

	require.Equal(t, Tags{}, secret.Tags)

	tag1 := "tag1"
	err = s.AddTag(ctx, secret.ID, tag1)
	require.NoError(t, err)

	loadedSecret, err = s.LoadSecretByName(ctx, secret.UserID, secret.Name)
	require.NoError(t, err)
	require.Equal(t, Tags{tag1}, loadedSecret.Tags)

	tag2 := "tag2"
	err = s.AddTag(ctx, secret.ID, tag2)
	require.NoError(t, err)

	loadedSecret, err = s.LoadSecretByName(ctx, secret.UserID, secret.Name)
	require.NoError(t, err)
	require.Equal(t, Tags{tag1, tag2}, loadedSecret.Tags)

	err = s.DeleteTag(ctx, secret.ID, tag1)
	require.NoError(t, err)

	loadedSecret, err = s.LoadSecretByName(ctx, secret.UserID, secret.Name)
	require.NoError(t, err)
	require.Equal(t, Tags{tag2}, loadedSecret.Tags)
}

func TestPgSQL_LoadSecrets_with_tags(t *testing.T) {
	var err error

	ctx := context.Background()
	s := setUp(ctx, t)

	user := createRandomUser(ctx, s, t)

	const tag1 = "13"
	const tag2 = "37"

	secretID := utils.NewUUID6()
	secret := &Secret{
		ID:     secretID,
		UserID: user.ID,
		Name:   "Card " + rand.RandomString(10),
		Tags:   Tags{},
		Kind:   api.KindBankCard,
		Value: &SecretBankCard{
			ID:     secretID,
			Name:   "KIRILL TITOV",
			Number: "1234 5678 9012 3456",
			Date:   "12/34/56",
			CVV:    "322",
		},
	}
	err = s.CreateSecret(ctx, secret)
	require.NoError(t, err)

	err = s.AddTag(ctx, secretID, tag1)
	require.NoError(t, err)
	err = s.AddTag(ctx, secretID, tag2)
	require.NoError(t, err)

	loadedSecrets, err := s.LoadSecrets(ctx, user.ID)
	require.NoError(t, err)
	loadedSecret := (loadedSecrets)[0]
	require.Len(t, loadedSecret.Tags, 2)
	require.Equal(t, loadedSecret.Tags, Tags{tag1, tag2})
}
