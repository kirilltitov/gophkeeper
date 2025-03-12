package app

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/container"
	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/storage"
	mockStorage "github.com/kirilltitov/gophkeeper/internal/storage/mocks"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func TestApplication_HandlerCreateSecretNote(t *testing.T) {
	a := Application{
		Gophkeeper: &gophkeeper.Gophkeeper{
			Config:    config.NewWithoutParsing(),
			Container: &container.Container{Storage: nil},
		},
	}

	userID := utils.NewUUID6()

	emptyStorage := func() storage.Storage {
		return mockStorage.NewMockStorage(t)
	}

	type input struct {
		body    string
		userID  *uuid.UUID
		storage func() storage.Storage
	}
	type want struct {
		code     int
		response string
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
				userID:  &userID,
				storage: emptyStorage,
			},
			want: want{
				code:     400,
				response: `{"success":false,"result":null,"error":"invalid input JSON"}`,
			},
		},
		{
			name: "Negative (invalid request 2)",
			input: input{
				body:    `{}`,
				userID:  &userID,
				storage: emptyStorage,
			},
			want: want{
				code:     400,
				response: `{"success":false,"result":null,"error":"invalid input JSON"}`,
			},
		},
		{
			name: "Negative (no auth)",
			input: input{
				body:    `{}`,
				storage: emptyStorage,
			},
			want: want{
				code: 401,
			},
		},
		{
			name: "Positive",
			input: input{
				body: `
					{
						"name": "secret note",
						"description": "some description",
						"value": {
							"body": "foo"
						}
					}
				`,
				userID: &userID,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						CreateSecret(mock.Anything, mock.Anything).
						Return(nil)
					return s
				},
			},
			want: want{
				code:     201,
				response: `{"success":true,"result":{"id": "<<PRESENCE>>"},"error":null}`,
			},
		},
		{
			name: "Negative (duplicate)",
			input: input{
				body: `
					{
						"name": "secret note",
						"description": "some description",
						"value": {
							"body": "foo"
						}
					}
				`,
				userID: &userID,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)
					s.
						EXPECT().
						CreateSecret(mock.Anything, mock.Anything).
						Return(storage.ErrDuplicateSecretFound)
					return s
				},
			},
			want: want{
				code:     409,
				response: `{"success":false,"result":null,"error":"secret with this name already exists"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.Gophkeeper.Container.Storage = tt.input.storage()

			r := httptest.NewRequest(
				http.MethodPost,
				"/api/secret/create/note",
				bytes.NewReader([]byte(tt.input.body)),
			)
			if tt.input.userID != nil {
				r = r.WithContext(utils.SetUserID(context.Background(), *tt.input.userID))
			}

			w := httptest.NewRecorder()

			a.HandlerCreateSecretNote(w, r)

			result := w.Result()
			defer result.Body.Close()

			actualResponse, err := io.ReadAll(result.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, result.StatusCode)

			if tt.want.response != "" {
				jsonassert.New(t).Assertf(string(actualResponse), tt.want.response)
			}
		})
	}
}
