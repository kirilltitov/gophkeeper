package app

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

type Claims struct {
	jwt.RegisteredClaims
}

type CtxUserIDKey struct{}

func (a *Application) WithAuthorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if userID := a.authorize(r); userID != nil {
			utils.Log.Infof("Authorized user %s by JWT cookie", userID.String())
			r = r.WithContext(context.WithValue(r.Context(), CtxUserIDKey{}, *userID))
		}

		next.ServeHTTP(w, r)
	}
}

func (a *Application) authorize(r *http.Request) *uuid.UUID {
	cfg := a.Gophkeeper.Config

	cookie, err := r.Cookie(cfg.JWTCookieName)
	if err != nil {
		utils.Log.WithError(err).Info("Could not authorize request")
		return nil
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		utils.Log.WithError(err).Info("Could not parse JWT")
		return nil
	}
	if !token.Valid {
		utils.Log.Info("JWT not valid")
		return nil
	}
	if claims.Subject == "" {
		utils.Log.Info("Missing Subject in JWT")
		return nil
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		utils.Log.WithError(err).Info("Invalid user ID UUID in JWT")
		return nil
	}

	return &userID
}

func (a *Application) getJWT(user storage.User, exp time.Time) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        utils.NewUUID6().String(),
				Subject:   user.ID.String(),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(exp),
			},
		},
	)

	return token.SignedString([]byte(a.Gophkeeper.Config.JWTSecret))
}

func (a *Application) CreateAuthCookie(user storage.User) (*http.Cookie, error) {
	cfg := a.Gophkeeper.Config

	exp := time.Now().Add(time.Second * time.Duration(cfg.JWTTimeToLive))
	token, err := a.getJWT(user, exp)
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:    cfg.JWTCookieName,
		Value:   token,
		Expires: exp,
	}, nil
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(CtxUserIDKey{}).(uuid.UUID)
	return userID, ok
}
