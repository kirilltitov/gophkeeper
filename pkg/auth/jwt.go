package auth

import "github.com/golang-jwt/jwt/v4"

// DefaultCookieName is a default JWT auth cookie name if not provided in env/CLI args.
const DefaultCookieName = "access_token"

// Claims contains all possible values stored in auth JWT.
type Claims struct {
	jwt.RegisteredClaims

	Login string `json:"login"` // Login is user's login.
}
