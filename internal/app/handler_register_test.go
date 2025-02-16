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
)

func TestApplication_HandlerRegister(t *testing.T) {
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
				body:    `{}`,
				storage: mockStorage.NewMockStorage(t),
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "Negative (invalid request 3)",
			input: input{
				body:    `{"login":"frankstrino","passworddd":"hesoyam"}`,
				storage: mockStorage.NewMockStorage(t),
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "Positive",
			input: input{
				body: `{"login":"frankstrino","password":"hesoyam"}`,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						CreateUser(mock.Anything, mock.Anything).
						Return(nil)
					return s
				}(),
			},
			want: want{
				code:      200,
				cookieSet: true,
			},
		},
		{
			name: "Negative (duplicate)",
			input: input{
				body: `{"login":"frankstrino","password":"hesoyam"}`,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						CreateUser(mock.Anything, mock.Anything).
						Return(storage.ErrDuplicateFound)
					return s
				}(),
			},
			want: want{
				code:      409,
				cookieSet: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.Gophkeeper.Container.Storage = tt.input.storage
			r := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader([]byte(tt.input.body)))
			w := httptest.NewRecorder()

			a.HandlerRegister(w, r)

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
