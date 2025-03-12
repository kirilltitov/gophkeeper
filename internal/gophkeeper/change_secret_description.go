package gophkeeper

import (
	"context"

	"github.com/google/uuid"
)

// ChangeSecretDescription changes a secret description.
func (g *Gophkeeper) ChangeSecretDescription(ctx context.Context, secretID uuid.UUID, description string) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	return g.Container.Storage.ChangeSecretDescription(ctx, secret.ID, description)
}
