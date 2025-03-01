package gophkeeper

import (
	"context"

	"github.com/google/uuid"
)

// AddTag adds a tag to an existing secret.
func (g *Gophkeeper) AddTag(ctx context.Context, secretID uuid.UUID, tag string) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	return g.Container.Storage.AddTag(ctx, secret.ID, tag)
}
