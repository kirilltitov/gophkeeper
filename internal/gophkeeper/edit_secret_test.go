package gophkeeper

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/kirilltitov/gophkeeper/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/container"
	"github.com/kirilltitov/gophkeeper/internal/storage"
	mockStorage "github.com/kirilltitov/gophkeeper/internal/storage/mocks"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func TestGophkeeper_EditSecretCredentials(t *testing.T) {
	cfg := config.NewWithoutParsing()
	cnt := container.Container{Storage: nil}

	g := New(cfg, &cnt)

	user := &storage.User{
		ID: utils.NewUUID6(),
	}
	secret := storage.Secret{
		ID:     utils.NewUUID6(),
		UserID: user.ID,
		Kind:   api.KindCredentials,
	}

	tests := []struct {
		name   string
		userID *uuid.UUID
		input  func() storage.Storage
		want   error
	}{
		{
			name:   "Positive",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&secret, nil)
				s.
					EXPECT().
					EditSecretCredentials(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
				return s
			},
			want: nil,
		},
		{
			name:   "Negative (wrong user)",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)

				wrongUserSecret := secret
				wrongUserSecret.UserID = utils.NewUUID6()
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&wrongUserSecret, nil)
				return s
			},
			want: ErrNoAuth,
		},
		{
			name:   "Negative (wrong kind)",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)

				wrongKindSecret := secret
				wrongKindSecret.Kind = api.KindBlob
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&wrongKindSecret, nil)
				return s
			},
			want: storage.ErrWrongKind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Container.Storage = tt.input()

			requestContext := context.Background()
			if tt.userID != nil {
				requestContext = utils.SetUserID(context.Background(), *tt.userID)
			}
			err := g.EditSecretCredentials(requestContext, secret.ID, "url", "login", "password")

			if tt.want != nil {
				assert.ErrorIs(t, err, tt.want)
			}
		})
	}
}

func TestGophkeeper_EditSecretBankCard(t *testing.T) {
	cfg := config.NewWithoutParsing()
	cnt := container.Container{Storage: nil}

	g := New(cfg, &cnt)

	user := &storage.User{
		ID: utils.NewUUID6(),
	}
	secret := storage.Secret{
		ID:     utils.NewUUID6(),
		UserID: user.ID,
		Kind:   api.KindBankCard,
	}

	tests := []struct {
		name   string
		userID *uuid.UUID
		input  func() storage.Storage
		want   error
	}{
		{
			name:   "Positive",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&secret, nil)
				s.
					EXPECT().
					EditSecretBankCard(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
				return s
			},
			want: nil,
		},
		{
			name:   "Negative (wrong user)",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)

				wrongUserSecret := secret
				wrongUserSecret.UserID = utils.NewUUID6()
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&wrongUserSecret, nil)
				return s
			},
			want: ErrNoAuth,
		},
		{
			name:   "Negative (wrong kind)",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)

				wrongKindSecret := secret
				wrongKindSecret.Kind = api.KindBlob
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&wrongKindSecret, nil)
				return s
			},
			want: storage.ErrWrongKind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Container.Storage = tt.input()

			requestContext := context.Background()
			if tt.userID != nil {
				requestContext = utils.SetUserID(context.Background(), *tt.userID)
			}
			err := g.EditSecretBankCard(
				requestContext,
				secret.ID,
				"name",
				"1234",
				"date",
				"cvv",
			)

			if tt.want != nil {
				assert.ErrorIs(t, err, tt.want)
			}
		})
	}
}

func TestGophkeeper_EditSecretBlob(t *testing.T) {
	cfg := config.NewWithoutParsing()
	cnt := container.Container{Storage: nil}

	g := New(cfg, &cnt)

	user := &storage.User{
		ID: utils.NewUUID6(),
	}
	secret := storage.Secret{
		ID:     utils.NewUUID6(),
		UserID: user.ID,
		Kind:   api.KindBlob,
	}

	tests := []struct {
		name   string
		userID *uuid.UUID
		input  func() storage.Storage
		want   error
	}{
		{
			name:   "Positive",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&secret, nil)
				s.
					EXPECT().
					EditSecretBlob(mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
				return s
			},
			want: nil,
		},
		{
			name:   "Negative (wrong user)",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)

				wrongUserSecret := secret
				wrongUserSecret.UserID = utils.NewUUID6()
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&wrongUserSecret, nil)
				return s
			},
			want: ErrNoAuth,
		},
		{
			name:   "Negative (wrong kind)",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)

				wrongKindSecret := secret
				wrongKindSecret.Kind = api.KindCredentials
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&wrongKindSecret, nil)
				return s
			},
			want: storage.ErrWrongKind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Container.Storage = tt.input()

			requestContext := context.Background()
			if tt.userID != nil {
				requestContext = utils.SetUserID(context.Background(), *tt.userID)
			}
			err := g.EditSecretBlob(
				requestContext,
				secret.ID,
				"name",
			)

			if tt.want != nil {
				assert.ErrorIs(t, err, tt.want)
			}
		})
	}
}

func TestGophkeeper_EditSecretNote(t *testing.T) {
	cfg := config.NewWithoutParsing()
	cnt := container.Container{Storage: nil}

	g := New(cfg, &cnt)

	user := &storage.User{
		ID: utils.NewUUID6(),
	}
	secret := storage.Secret{
		ID:     utils.NewUUID6(),
		UserID: user.ID,
		Kind:   api.KindNote,
	}

	tests := []struct {
		name   string
		userID *uuid.UUID
		input  func() storage.Storage
		want   error
	}{
		{
			name:   "Positive",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&secret, nil)
				s.
					EXPECT().
					EditSecretNote(mock.Anything, mock.Anything, mock.Anything).
					Return(nil)
				return s
			},
			want: nil,
		},
		{
			name:   "Negative (wrong user)",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)

				wrongUserSecret := secret
				wrongUserSecret.UserID = utils.NewUUID6()
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&wrongUserSecret, nil)
				return s
			},
			want: ErrNoAuth,
		},
		{
			name:   "Negative (wrong kind)",
			userID: &user.ID,
			input: func() storage.Storage {
				s := mockStorage.NewMockStorage(t)

				wrongKindSecret := secret
				wrongKindSecret.Kind = api.KindCredentials
				s.
					EXPECT().
					LoadSecretByID(mock.Anything, mock.Anything).
					Return(&wrongKindSecret, nil)
				return s
			},
			want: storage.ErrWrongKind,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Container.Storage = tt.input()

			requestContext := context.Background()
			if tt.userID != nil {
				requestContext = utils.SetUserID(context.Background(), *tt.userID)
			}
			err := g.EditSecretNote(
				requestContext,
				secret.ID,
				"name",
			)

			if tt.want != nil {
				assert.ErrorIs(t, err, tt.want)
			}
		})
	}
}
