package gophkeeper

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/container"
	"github.com/kirilltitov/gophkeeper/internal/storage"
	mockStorage "github.com/kirilltitov/gophkeeper/internal/storage/mocks"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func TestGophkeeper_EditSecretDescription(t *testing.T) {
	cfg := config.NewWithoutParsing()
	cnt := container.Container{Storage: nil}

	g := New(cfg, &cnt)

	user := &storage.User{
		ID: utils.NewUUID6(),
	}
	secret := storage.Secret{
		ID:     utils.NewUUID6(),
		UserID: user.ID,
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
					ChangeSecretDescription(mock.Anything, mock.Anything, mock.Anything).
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Container.Storage = tt.input()

			requestContext := context.Background()
			if tt.userID != nil {
				requestContext = utils.SetUserID(context.Background(), *tt.userID)
			}
			err := g.ChangeSecretDescription(requestContext, secret.ID, "foo")

			if tt.want != nil {
				assert.ErrorIs(t, tt.want, err)
			}
		})
	}
}
