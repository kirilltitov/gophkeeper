package gophkeeper

import (
	"context"

	"github.com/google/uuid"
)

func (g *Gophkeeper) RenameSecret(ctx context.Context, secretID uuid.UUID, name string) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	return g.Container.Storage.RenameSecret(ctx, secret.ID, name)
}
