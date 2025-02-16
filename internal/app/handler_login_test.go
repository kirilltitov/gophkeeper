package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/container"
	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/storage"
	mockStorage "github.com/kirilltitov/gophkeeper/internal/storage/mocks"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func TestApplication_HandlerLogin(t *testing.T) {
	t.Parallel()

	a := Application{
		Gophkeeper: gophkeeper.Gophkeeper{
			Config:    config.NewWithoutParsing(),
			Container: &container.Container{Storage: nil},
		},
	}

	type input struct {
		body    string
		storage storage.Storage
	}
	type want struct {
		code      int
		cookieSet bool
	}
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name: "Negative (invalid request 1)",
			input: input{
				body:    `invalid`,
				storage: mockStorage.NewMockStorage(t),
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "Negative (invalid request 2)",
			input: input{
				body: `{}`,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						LoadUser(mock.Anything, mock.Anything).
						Return(nil, gophkeeper.ErrAuthFailed)
					return s
				}(),
			},
			want: want{
				code: 401,
			},
		},
		{
			name: "Negative (invalid request 3)",
			input: input{
				body: `{"login":"frankstrino","passworddd":"hesoyam"}`,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						LoadUser(mock.Anything, mock.Anything).
						Return(nil, gophkeeper.ErrAuthFailed)
					return s
				}(),
			},
			want: want{
				code: 401,
			},
		},
		{
			name: "Positive",
			input: func() input {
				userID := utils.NewUUID6()
				user := storage.NewUser(userID, "frankstrino", "hesoyam")

				return input{
					body: `{"login":"frankstrino","password":"hesoyam"}`,
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
				code:      200,
				cookieSet: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.Gophkeeper.Container.Storage = tt.input.storage
			r := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewReader([]byte(tt.input.body)))
			w := httptest.NewRecorder()

			a.HandlerLogin(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			if tt.want.cookieSet {
				assert.NotEmpty(t, result.Cookies())
			} else {
				assert.Empty(t, result.Cookies())
			}
		})
	}
}
