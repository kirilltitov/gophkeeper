package gophkeeper

import (
	"context"

	"github.com/google/uuid"
)

func (g *Gophkeeper) DeleteTag(ctx context.Context, secretID uuid.UUID, tag string) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	return g.Container.Storage.DeleteTag(ctx, secret.ID, tag)
}
