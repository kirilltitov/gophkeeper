package app

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/pkg/auth"
)

// WithAuthorization is a middleware for an HTTP server authorizing user with a JWT cookie.
// If successful, user ID is set to Context under [utils.CtxUserIDKey] key.
// This user ID must be trusted.
func (a *Application) WithAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if userID := a.authorize(r); userID != nil {
			utils.Log.Infof("Authorized user %s by JWT cookie", userID.String())
			r = r.WithContext(utils.SetUserID(r.Context(), *userID))
		}

		next.ServeHTTP(w, r)
	})
}

// CreateAuthCookie creates and returns a new authorization cookie with a signed JWT.
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

func (a *Application) authorize(r *http.Request) *uuid.UUID {
	cfg := a.Gophkeeper.Config

	cookie, err := r.Cookie(cfg.JWTCookieName)
	if err != nil {
		utils.Log.WithError(err).Info("Could not authorize request")
		return nil
	}

	claims := &auth.Claims{}

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
		auth.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        utils.NewUUID6().String(),
				Subject:   user.ID.String(),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(exp),
			},
			Login: user.Login,
		},
	)

	return token.SignedString([]byte(a.Gophkeeper.Config.JWTSecret))
}
