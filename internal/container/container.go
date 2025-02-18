package container

import (
	"context"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/storage"
)

type Container struct {
	Storage storage.Storage
}

func New(ctx context.Context, cfg *config.Config) (*Container, error) {
	s, err := newPgSQLStorage(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Container{
		Storage: *s,
	}, nil
}
