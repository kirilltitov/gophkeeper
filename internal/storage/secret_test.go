package storage

import (
	"context"
	"testing"

	"github.com/kirilltitov/gophkeeper/pkg/api"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/internal/utils/rand"
)

func createRandomSecret(t *testing.T, ctx context.Context, s Storage) *Secret {
	return createRandomSecretForUser(t, ctx, s, createRandomUser(ctx, s, t))
}

func createRandomSecretForUser(t *testing.T, ctx context.Context, s Storage, user *User) *Secret {
	var err error

	secretID := utils.NewUUID6()
	secretName := "Card " + rand.RandomString(10)
	secret := &Secret{
		ID:     secretID,
		UserID: user.ID,
		Name:   secretName,
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

	return secret
}

func TestPgSQL_CreateSecret(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	var err error

	secret := createRandomSecret(t, ctx, s)

	secret2ID := utils.NewUUID6()
	secret2SameName := &Secret{
		ID:     secret2ID,
		UserID: secret.UserID,
		Name:   secret.Name,
		Kind:   api.KindNote,
		Value: &SecretNote{
			ID:   secret2ID,
			Body: "secret",
		},
	}
	err = s.CreateSecret(ctx, secret2SameName)
	require.ErrorIs(t, err, ErrDuplicateSecretFound)

	secret3SameID := &Secret{
		ID:     secret.ID,
		UserID: secret.UserID,
		Name:   "Note " + rand.RandomString(10),
		Kind:   api.KindNote,
		Value: &SecretNote{
			ID:   secret.ID,
			Body: "secret",
		},
	}
	err = s.CreateSecret(ctx, secret3SameID)
	require.ErrorIs(t, err, ErrDuplicateSecretFound)
}

func TestPgSQL_DeleteSecret(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	var err error

	secret := createRandomSecret(t, ctx, s)

	err = s.DeleteSecret(ctx, secret.ID)
	require.NoError(t, err)

	_, err = s.LoadSecretByID(ctx, secret.ID)
	require.ErrorIs(t, err, ErrNotFound)
}

func TestPgSQL_RenameSecret(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	var err error

	user := createRandomUser(ctx, s, t)

	existingSecret := createRandomSecretForUser(t, ctx, s, user)
	newSecret := createRandomSecretForUser(t, ctx, s, user)

	err = s.RenameSecret(ctx, newSecret.ID, existingSecret.Name)
	require.ErrorIs(t, err, ErrDuplicateSecretFound)

	err = s.RenameSecret(ctx, newSecret.ID, rand.RandomString(10))
	require.NoError(t, err)
}

func TestPgSQL_LoadSecretByName(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	loadedSecret, err := s.LoadSecretByName(ctx, utils.NewUUID6(), rand.RandomString(10))
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, loadedSecret)

	secret := createRandomSecret(t, ctx, s)

	loadedSecret, err = s.LoadSecretByName(ctx, secret.UserID, secret.Name)
	require.NoError(t, err)
	require.NotNil(t, loadedSecret)
	require.Equal(t, secret, loadedSecret)
}

func TestPgSQL_LoadSecretByID(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	loadedSecret, err := s.LoadSecretByID(ctx, utils.NewUUID6())
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, loadedSecret)

	secret := createRandomSecret(t, ctx, s)

	loadedSecret, err = s.LoadSecretByID(ctx, secret.ID)
	require.NoError(t, err)
	require.NotNil(t, loadedSecret)
	require.Equal(t, secret, loadedSecret)
}

func TestPgSQL_LoadSecrets(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	var err error
	user := createRandomUser(ctx, s, t)

	loadedSecrets, err := s.LoadSecrets(ctx, user.ID)
	require.NoError(t, err)
	require.Len(t, *loadedSecrets, 0)

	const numSecrets = 5
	for i := 0; i < numSecrets; i++ {
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
	}

	loadedSecrets, err = s.LoadSecrets(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, loadedSecrets)
	require.Len(t, *loadedSecrets, numSecrets)
}

func TestPgSQL_LoadSecret_values(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	user := createRandomUser(ctx, s, t)

	type input struct {
		name  string
		input *Secret
		want  *Secret
	}
	tests := []input{
		func() input {
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
			return input{
				name:  "Positive: Bank Card",
				input: secret,
				want:  secret,
			}
		}(),
		func() input {
			secretID := utils.NewUUID6()
			secret := &Secret{
				ID:     secretID,
				UserID: user.ID,
				Name:   "Blob " + rand.RandomString(10),
				Tags:   Tags{},
				Kind:   api.KindBlob,
				Value: &SecretBlob{
					ID:   secretID,
					Body: rand.RandomString(16),
				},
			}
			return input{
				name:  "Positive: Blob",
				input: secret,
				want:  secret,
			}
		}(),
		func() input {
			secretID := utils.NewUUID6()
			secret := &Secret{
				ID:     secretID,
				UserID: user.ID,
				Name:   "Note " + rand.RandomString(10),
				Tags:   Tags{},
				Kind:   api.KindNote,
				Value: &SecretNote{
					ID:   secretID,
					Body: rand.RandomString(16),
				},
			}
			return input{
				name:  "Positive: Note",
				input: secret,
				want:  secret,
			}
		}(),
		func() input {
			secretID := utils.NewUUID6()
			secret := &Secret{
				ID:     secretID,
				UserID: user.ID,
				Name:   "Credentials " + rand.RandomString(10),
				Tags:   Tags{},
				Kind:   api.KindCredentials,
				Value: &SecretCredentials{
					ID:       secretID,
					Login:    rand.RandomString(10),
					Password: rand.RandomString(10),
				},
			}
			return input{
				name:  "Positive: Credentials",
				input: secret,
				want:  secret,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.CreateSecret(ctx, tt.input)
			require.NoError(t, err)

			loadedSecret, err := s.LoadSecretByName(ctx, user.ID, tt.input.Name)
			require.NoError(t, err)
			require.NotNil(t, loadedSecret)
			require.Equal(t, tt.input, loadedSecret)
		})
	}
}

func TestPgSQL_CreateSecret_kind_mismatch(t *testing.T) {
	ctx := context.Background()
	s := setUp(ctx, t)

	user := createRandomUser(ctx, s, t)

	type input struct {
		name  string
		input *Secret
	}
	tests := []input{
		func() input {
			secretID := utils.NewUUID6()
			secret := &Secret{
				ID:     secretID,
				UserID: user.ID,
				Name:   "Card " + rand.RandomString(10),
				Kind:   api.KindNote,
				Value:  &SecretBankCard{},
			}
			return input{
				name:  "Negative: note and card",
				input: secret,
			}
		}(),
		func() input {
			secretID := utils.NewUUID6()
			secret := &Secret{
				ID:     secretID,
				UserID: user.ID,
				Name:   "Blob " + rand.RandomString(10),
				Kind:   api.KindBankCard,
				Value:  &SecretBlob{},
			}
			return input{
				name:  "Negative: card and blob",
				input: secret,
			}
		}(),
		func() input {
			secretID := utils.NewUUID6()
			secret := &Secret{
				ID:     secretID,
				UserID: user.ID,
				Name:   "Note " + rand.RandomString(10),
				Kind:   api.KindBlob,
				Value: &SecretNote{
					ID:   secretID,
					Body: rand.RandomString(16),
				},
			}
			return input{
				name:  "Negative: blob and note",
				input: secret,
			}
		}(),
		func() input {
			secretID := utils.NewUUID6()
			secret := &Secret{
				ID:     secretID,
				UserID: user.ID,
				Name:   "Credentials " + rand.RandomString(10),
				Kind:   api.KindCredentials,
				Value:  &SecretNote{},
			}
			return input{
				name:  "Negative: credentials and note",
				input: secret,
			}
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.CreateSecret(ctx, tt.input)
			require.ErrorIs(t, err, ErrWrongKind)
		})
	}
}
