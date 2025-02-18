package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func TestApplication_WithAuthorization(t *testing.T) {
	a := Application{
		Gophkeeper: &gophkeeper.Gophkeeper{
			Config: config.NewWithoutParsing(),
		},
	}

	validUser := storage.User{ID: utils.NewUUID6()}

	tests := []struct {
		name   string
		cookie *http.Cookie
		want   *uuid.UUID
	}{
		{
			name: "Positive",
			cookie: &http.Cookie{
				Name: "access_token",
				Value: func() string {
					token, _ := a.getJWT(validUser, time.Now().Add(time.Second*10))
					return token
				}(),
			},
			want: &validUser.ID,
		},
		{
			name: "Negative (no cookie)",
		},
		{
			name: "Negative (not a JWT)",
			cookie: &http.Cookie{
				Name:  "access_token",
				Value: "NOT A JWT",
			},
		},
		{
			name: "Negative (invalid signature)",
			cookie: &http.Cookie{
				Name: "access_token",
				Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIwMWVmZWNhNS0yOTllLTY4YjQtODFiYi05YWQ1ODQ0YzNmY" +
					"jgiLCJleHAiOjE3Mzk3Mzc5NjgsImlhdCI6MTczOTczNzk1OCwianRpIjoiMDFlZmVjYTUtMjk5ZS02YmEyLTgxYmItOWF" +
					"kNTg0NGMzZmI4In0." +
					"INVALID SIGNATURE",
			},
		},
		{
			name: "Negative (valid expired JWT)",
			cookie: &http.Cookie{
				Name: "access_token",
				Value: func() string {
					token, _ := a.getJWT(validUser, time.Now().Add(-time.Second))
					return token
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", http.NoBody)
			require.NoError(t, err)

			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			handler := a.WithAuthorization(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				userID, ok := GetUserID(r.Context())

				if tt.want != nil {
					require.Equal(t, tt.want, &userID)
					require.True(t, ok)
				} else {
					require.False(t, ok)
				}
			}))
			handler.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}
