package gophkeeper

import (
	"context"

	"github.com/google/uuid"
)

// DeleteSecret deletes an existing secret.
func (g *Gophkeeper) DeleteSecret(ctx context.Context, secretID uuid.UUID) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	return g.Container.Storage.DeleteSecret(ctx, secret.ID)
}
