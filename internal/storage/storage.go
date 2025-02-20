package storage

import (
	"context"

	"github.com/google/uuid"
)

//go:generate mockery
type Storage interface {
	// CreateUser создает нового пользователя в хранилище.
	CreateUser(ctx context.Context, user User) error

	// LoadUser пытается найти пользователя в хранилище по заданному логину.
	LoadUser(ctx context.Context, login string) (*User, error)

	CreateSecret(ctx context.Context, secret *Secret) error

	LoadSecret(ctx context.Context, userID uuid.UUID, name string) (*Secret, error)

	LoadSecrets(ctx context.Context, user User) (*[]Secret, error)

	AddTag(ctx context.Context, secretID uuid.UUID, tag string) error

	DeleteTag(ctx context.Context, secretID uuid.UUID, tag string) error

	// Close закрывает соединение с хранилищем.
	Close()
}
