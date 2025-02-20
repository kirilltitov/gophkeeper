package gophkeeper

import (
	"context"

	"github.com/kirilltitov/gophkeeper/internal/storage"
)

func (g *Gophkeeper) CreateSecret(ctx context.Context, secret *storage.Secret) error {
	return g.Container.Storage.CreateSecret(ctx, secret)
}
