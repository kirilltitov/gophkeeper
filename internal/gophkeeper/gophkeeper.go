package gophkeeper

import (
	"context"

	"github.com/google/uuid"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/container"
	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// Gophkeeper является объектом, инкапсулирующим в себе бизнес-логику сервиса по хранению секретов.
type Gophkeeper struct {
	Config    *config.Config
	Container *container.Container
}

// New создает, конфигурирует и возвращает экземпляр объекта сервиса.
func New(cfg *config.Config, cnt *container.Container) *Gophkeeper {
	return &Gophkeeper{Config: cfg, Container: cnt}
}

func (g *Gophkeeper) loadSecretAndAuthorize(ctx context.Context, secretID uuid.UUID) (*storage.Secret, error) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		return nil, ErrNoAuth
	}

	secret, err := g.Container.Storage.LoadSecretByID(ctx, secretID)
	if err != nil {
		return nil, err
	}

	if secret.UserID != userID {
		return nil, ErrNoAuth
	}

	return secret, nil
}
