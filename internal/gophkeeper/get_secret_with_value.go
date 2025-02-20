package gophkeeper

import (
	"context"

	"github.com/google/uuid"

	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func (g *Gophkeeper) GetSecretWithValueByName(
	ctx context.Context,
	name string,
) (*storage.Secret, error) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		return nil, ErrNoAuth
	}

	return g.Container.Storage.LoadSecretByName(ctx, userID, name)
}

func (g *Gophkeeper) GetSecretWithValueByID(
	ctx context.Context,
	secretID uuid.UUID,
) (*storage.Secret, error) {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
