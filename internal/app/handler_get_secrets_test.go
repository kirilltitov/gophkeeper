package app

import (
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
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

func TestApplication_HandlerGetSecrets(t *testing.T) {
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
			name: "Negative (no auth)",
			input: input{
				storage: emptyStorage,
			},
			want: want{
				code:     401,
				response: `{"success":false,"result":null,"error":"unauthorized"}`,
			},
		},
		{
			name: "Positive",
			input: input{
				userID: &userID,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)

					secret1ID := utils.NewUUID6()
					secret2ID := utils.NewUUID6()
					result := []*storage.Secret{
						{
							ID:          secret1ID,
							UserID:      userID,
							Name:        "foo",
							Description: "foo description",
							Tags:        storage.Tags{"bar", "baz"},
							Kind:        api.KindNote,
							IsEncrypted: false,
							Value: &storage.SecretNote{
								ID:   secret1ID,
								Body: "foo body",
							},
						},
						{
							ID:          secret2ID,
							UserID:      userID,
							Name:        "bar",
							Description: "bar description",
							Tags:        storage.Tags{},
							Kind:        api.KindCredentials,
							IsEncrypted: false,
							Value: &storage.SecretCredentials{
								ID:       secret2ID,
								URL:      "someurl",
								Login:    "teonoman",
								Password: "megapass",
							},
						},
					}

					s.
						EXPECT().
						LoadSecrets(mock.Anything, mock.Anything).
						Return(result, nil)
					return s
				},
			},
			want: want{
				code: 200,
				response: `
					{
						"success": true,
						"result": [
							{
								"id": "<<PRESENCE>>",
								"user_id": "<<PRESENCE>>",
								"name": "foo",
								"description": "foo description",
								"tags": ["bar","baz"],
								"kind": "note",
								"is_encrypted": false,
								"value": {
									"id": "<<PRESENCE>>",
									"body": "foo body"
								}
							},
							{
								"id": "<<PRESENCE>>",
								"user_id": "<<PRESENCE>>",
								"name": "bar",
								"description": "bar description",
								"tags": [],
								"kind": "credentials",
								"is_encrypted": false,
								"value": {
									"id": "<<PRESENCE>>",
									"url": "someurl",
									"login": "teonoman",
									"password": "megapass"
								}
							}
						],
						"error": null
					}
				`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a.Gophkeeper.Container.Storage = tt.input.storage()

			r := httptest.NewRequest(
				http.MethodGet,
				"/api/secret/list",
				http.NoBody,
			)
			if tt.input.userID != nil {
				r = r.WithContext(utils.SetUserID(context.Background(), *tt.input.userID))
			}

			w := httptest.NewRecorder()

			a.HandlerGetSecrets(w, r)

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
