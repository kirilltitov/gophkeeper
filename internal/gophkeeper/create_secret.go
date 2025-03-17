package gophkeeper

import (
	"context"

	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// CreateSecret creates a new secret.
func (g *Gophkeeper) CreateSecret(ctx context.Context, secret *storage.Secret) error {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		return ErrNoAuth
	}

	secretID := utils.NewUUID6()
	secret.ID = secretID
	secret.Value.SetID(secretID)
	secret.Kind = secret.Value.Kind()
	secret.UserID = userID

	return g.Container.Storage.CreateSecret(ctx, secret)
}
