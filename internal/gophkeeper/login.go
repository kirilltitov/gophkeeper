package gophkeeper

import (
	"context"

	"github.com/kirilltitov/gophkeeper/internal/storage"
)

// Login authenticates a user with given login and password.
func (g *Gophkeeper) Login(ctx context.Context, login string, password string) (*storage.User, error) {
	user, err := g.Container.Storage.LoadUser(ctx, login)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.IsValidPassword(password) {
		return nil, ErrAuthFailed
	}

	return user, nil
}
