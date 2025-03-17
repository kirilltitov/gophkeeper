package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"

	"github.com/kirilltitov/gophkeeper/pkg/auth"
)

const (
	authFile = "auth.jwt"
)

func getJWTFileName() string {
	return fmt.Sprintf("%s/%s", getConfigDir(), authFile)
}

func storeJWT(jwt string) error {
	if err := os.WriteFile(getJWTFileName(), []byte(jwt), 0o660); err != nil {
		return err
	}

	return nil
}

func getAuthJWTFileContents() (string, error) {
	bytes, err := os.ReadFile(getJWTFileName())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}

	return string(bytes), nil
}

func getAuthClaimsFromString(jwtString string) (*auth.Claims, error) {
	claims := &auth.Claims{}

	if jwtString == "" {
		return claims, nil
	}

	_, err := jwt.ParseWithClaims(jwtString, claims, nil)

	if err != nil && !errors.Is(err, jwt.ErrTokenUnverifiable) {
		return nil, err
	}

	return claims, nil
}

func getAuthClaims() (*auth.Claims, error) {
	jwtString, err := getAuthJWTFileContents()
	if err != nil {
		return nil, err
	}

	return getAuthClaimsFromString(jwtString)
}

func authenticate() (string, error) {
	jwtString, err := getAuthJWTFileContents()
	if err != nil {
		return "", err
	}

	claims, err := getAuthClaimsFromString(jwtString)
	if err != nil {
		return "", err
	}

	if claims.Login == "" {
		return "", errNoAuth
	}

	if err := claims.Valid(); err != nil {
		return "", errAuthExpired
	}

	return jwtString, nil
}

func findAuthCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	return nil
}
