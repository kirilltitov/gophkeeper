package gophkeeper

import (
	"context"

	"github.com/google/uuid"
)

// RenameSecret renames a secret.
func (g *Gophkeeper) RenameSecret(ctx context.Context, secretID uuid.UUID, name string) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	return g.Container.Storage.RenameSecret(ctx, secret.ID, name)
}
