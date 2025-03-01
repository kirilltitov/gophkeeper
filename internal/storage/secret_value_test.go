package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/internal/utils/rand"
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

func TestPgSQL_EditSecretBlob(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	var err error
	user := createRandomUser(ctx, s, t)

	secretID := utils.NewUUID6()
	secret := &Secret{
		ID:     secretID,
		UserID: user.ID,
		Name:   "Blob " + rand.RandomString(10),
		Kind:   api.KindBlob,
		Value: &SecretBlob{
			ID:   secretID,
			Body: "someblob",
		},
	}
	err = s.CreateSecret(ctx, secret)
	require.NoError(t, err)

	newBody := "somenewblob"
	err = s.EditSecretBlob(ctx, secret, newBody)
	require.NoError(t, err)

	loadedSecret, err := s.LoadSecretByName(ctx, secret.UserID, secret.Name)
	require.NoError(t, err)
	require.NotNil(t, loadedSecret)
	require.Equal(t, newBody, loadedSecret.Value.(*SecretBlob).Body)
}

func TestPgSQL_EditSecretNote(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	var err error
	user := createRandomUser(ctx, s, t)

	secretID := utils.NewUUID6()
	secret := &Secret{
		ID:     secretID,
		UserID: user.ID,
		Name:   "Note " + rand.RandomString(10),
		Kind:   api.KindNote,
		Value: &SecretNote{
			ID:   secretID,
			Body: "somenote",
		},
	}
	err = s.CreateSecret(ctx, secret)
	require.NoError(t, err)

	newBody := "somenewnote"
	err = s.EditSecretNote(ctx, secret, newBody)
	require.NoError(t, err)

	loadedSecret, err := s.LoadSecretByName(ctx, secret.UserID, secret.Name)
	require.NoError(t, err)
	require.NotNil(t, loadedSecret)
	require.Equal(t, newBody, loadedSecret.Value.(*SecretNote).Body)
}

func TestPgSQL_EditSecretBankCard(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	var err error
	user := createRandomUser(ctx, s, t)

	secretID := utils.NewUUID6()
	secret := &Secret{
		ID:     secretID,
		UserID: user.ID,
		Name:   "Card " + rand.RandomString(10),
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

	newCard := &SecretBankCard{
		ID:     secretID,
		Name:   "KIRILLIUS TITOV",
		Number: "1111 2222 3333 4444",
		Date:   "09/03/1989",
		CVV:    "1337",
	}
	err = s.EditSecretBankCard(ctx, secret, newCard.Name, newCard.Number, newCard.Date, newCard.CVV)
	require.NoError(t, err)

	loadedSecret, err := s.LoadSecretByName(ctx, secret.UserID, secret.Name)
	require.NoError(t, err)
	require.NotNil(t, loadedSecret)
	require.Equal(t, newCard, loadedSecret.Value.(*SecretBankCard))
}

func TestPgSQL_EditSecretCredentials(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	var err error
	user := createRandomUser(ctx, s, t)

	secretID := utils.NewUUID6()
	secret := &Secret{
		ID:     secretID,
		UserID: user.ID,
		Name:   "Card " + rand.RandomString(10),
		Kind:   api.KindCredentials,
		Value: &SecretCredentials{
			ID:       secretID,
			Login:    "teonoman",
			Password: "megapass",
		},
	}
	err = s.CreateSecret(ctx, secret)
	require.NoError(t, err)

	newCredentials := &SecretCredentials{
		ID:       secretID,
		Login:    "teonoman2",
		Password: "megapass2",
	}
	err = s.EditSecretCredentials(ctx, secret, newCredentials.Login, newCredentials.Password)
	require.NoError(t, err)

	loadedSecret, err := s.LoadSecretByName(ctx, secret.UserID, secret.Name)
	require.NoError(t, err)
	require.NotNil(t, loadedSecret)
	require.Equal(t, newCredentials, loadedSecret.Value.(*SecretCredentials))
}
