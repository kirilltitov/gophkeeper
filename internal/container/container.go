package container

import (
	"context"

	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/storage"
)

// Container is a dependency container.
type Container struct {
	Storage storage.Storage // Storage is an interface to storage.
}

// New creates and returns a fully configured container.
func New(ctx context.Context, cfg *config.Config) (*Container, error) {
	s, err := newPgSQLStorage(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Container{
		Storage: s,
	}, nil
}
