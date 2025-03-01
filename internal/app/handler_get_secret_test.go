package app

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
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

func TestApplication_HandlerGetSecret(t *testing.T) {
	a := Application{
		Gophkeeper: &gophkeeper.Gophkeeper{
			Config:    config.NewWithoutParsing(),
			Container: &container.Container{Storage: nil},
		},
	}

	userID := utils.NewUUID6()
	secretID := utils.NewUUID6()

	emptyStorage := func() storage.Storage {
		return mockStorage.NewMockStorage(t)
	}

	type input struct {
		secretID string
		userID   *uuid.UUID
		storage  func() storage.Storage
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
				secretID: utils.NewUUID6().String(),
				storage:  emptyStorage,
			},
			want: want{
				code: 401,
			},
		},
		{
			name: "Negative (no ID)",
			input: input{
				userID:  &userID,
				storage: emptyStorage,
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "Positive",
			input: input{
				secretID: secretID.String(),
				userID:   &userID,
				storage: func() storage.Storage {
					s := mockStorage.NewMockStorage(t)

					secret := &storage.Secret{
						ID:          secretID,
						UserID:      userID,
						Name:        "my secret card",
						Tags:        storage.Tags{"MIR"},
						Kind:        api.KindBankCard,
						IsEncrypted: true,
						Value: &storage.SecretBankCard{
							ID:     secretID,
							Name:   "KIRILL TITOV",
							Number: "1234 5678 9012 3456",
							Date:   "12/34/5678",
							CVV:    "322",
						},
					}

					s.
						EXPECT().
						LoadSecretByID(mock.Anything, mock.Anything).
						Return(secret, nil)
					return s
				},
			},
			want: want{
				code: 200,
				response: `
					{
						"success": true,
						"result": {
							"id": "` + secretID.String() + `",
							"user_id": "` + userID.String() + `",
							"name": "my secret card",
							"tags": ["MIR"],
							"kind": "bank_card",
							"is_encrypted": true,
							"value": {
								"id": "` + secretID.String() + `",
								"name": "KIRILL TITOV",
								"number": "1234 5678 9012 3456",
								"date": "12/34/5678",
								"cvv": "322"
							}
						},
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
				"/api/secret/get/2a9186b1-d39f-49cb-99a9-b6e8a25293a2",
				http.NoBody,
			)
			if tt.input.userID != nil {
				r = r.WithContext(utils.SetUserID(context.Background(), *tt.input.userID))
			}

			if tt.input.secretID != "" {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("ID", tt.input.secretID)
				r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			}

			w := httptest.NewRecorder()

			a.HandlerGetSecret(w, r)

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
