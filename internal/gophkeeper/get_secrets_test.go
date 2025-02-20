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

func TestGophkeeper_GetSecrets(t *testing.T) {
	cfg := config.NewWithoutParsing()
	cnt := container.Container{Storage: nil}

	g := New(cfg, &cnt)

	user := &storage.User{
		ID: utils.NewUUID6(),
	}
	secret := storage.Secret{
		ID:     utils.NewUUID6(),
		UserID: user.ID,
		Kind:   storage.KindNote,
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
					LoadSecrets(mock.Anything, mock.Anything).
					Return(&[]storage.Secret{secret}, nil)
				return s
			},
			want: nil,
		},
		{
			name: "Negative (no user)",
			input: func() storage.Storage {
				return mockStorage.NewMockStorage(t)
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
			_, err := g.GetSecrets(
				requestContext,
			)

			if tt.want != nil {
				assert.ErrorIs(t, err, tt.want)
			}
		})
	}
}
