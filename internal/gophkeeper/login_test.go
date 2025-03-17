package gophkeeper

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/container"
	"github.com/kirilltitov/gophkeeper/internal/storage"
	mockStorage "github.com/kirilltitov/gophkeeper/internal/storage/mocks"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func TestGophermart_Login(t *testing.T) {
	cfg := config.NewWithoutParsing()
	ctx := context.Background()
	cnt := container.Container{Storage: nil}

	g := New(cfg, &cnt)

	type input struct {
		login    string
		password string
		storage  storage.Storage
	}
	type want struct {
		err  error
		user *storage.User
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "Negative (wrong password)",
			input: input{
				login:    `some`,
				password: `incorrect`,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						LoadUser(mock.Anything, mock.Anything).
						Return(
							&storage.User{
								ID:        uuid.UUID{},
								Login:     "foo",
								Password:  "bar",
								CreatedAt: time.Now(),
							},
							nil,
						)
					return s
				}(),
			},
			want: want{
				err: ErrAuthFailed,
			},
		},
		{
			name: "Positive",
			input: func() input {
				userID := utils.NewUUID6()
				user := storage.NewUser(userID, "frankstrino", "hesoyam")

				return input{
					login:    "frankstrino",
					password: "hesoyam",
					storage: func() storage.Storage {
						s := mockStorage.NewMockStorage(t)
						s.
							EXPECT().
							LoadUser(mock.Anything, mock.Anything).
							Return(&user, nil)
						return s
					}(),
				}
			}(),
			want: want{
				user: &storage.User{
					Login:    "frankstrino",
					Password: "hesoyam",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g.Container.Storage = tt.input.storage

			user, err := g.Login(ctx, tt.input.login, tt.input.password)

			if tt.want.err != nil {
				assert.ErrorIs(t, tt.want.err, err)
			}
			if tt.want.user != nil {
				require.NotNil(t, user)
				assert.Equal(t, tt.want.user.Login, user.Login)
			}
		})
	}
}
