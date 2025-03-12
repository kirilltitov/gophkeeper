package gophkeeper

import (
	"context"

	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// GetSecrets returns all secrets for current user.
func (g *Gophkeeper) GetSecrets(ctx context.Context) ([]*storage.Secret, error) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		return nil, ErrNoAuth
	}

	return g.Container.Storage.LoadSecrets(ctx, userID)
}
