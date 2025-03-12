package gophkeeper

import (
	"context"

	"github.com/google/uuid"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/container"
	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// Gophkeeper is object encapsulating all business logic of Gophkeeper service.
type Gophkeeper struct {
	Config    *config.Config       // Config is service configuration.
	Container *container.Container // Container contains all service dependencies.
}

// New creates and returns a new instance of Gophkeeper instance.
func New(cfg *config.Config, cnt *container.Container) *Gophkeeper {
	return &Gophkeeper{Config: cfg, Container: cnt}
}

func (g *Gophkeeper) loadSecretAndAuthorize(ctx context.Context, secretID uuid.UUID) (*storage.Secret, error) {
	userID, _ := utils.GetUserID(ctx)

	secret, err := g.Container.Storage.LoadSecretByID(ctx, secretID)
	if err != nil {
		return nil, err
	}

	if secret.UserID != userID {
		return nil, ErrNoAuth
	}

	return secret, nil
}
