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

	RenameSecret(ctx context.Context, secretID uuid.UUID, name string) error

	DeleteSecret(ctx context.Context, secretID uuid.UUID) error

	EditSecretCredentials(ctx context.Context, secret *Secret, login string, password string) error

	EditSecretNote(ctx context.Context, secret *Secret, body string) error

	EditSecretBlob(ctx context.Context, secret *Secret, body string) error

	EditSecretBankCard(ctx context.Context, secret *Secret, name, number, date, cvv string) error

	LoadSecretByName(ctx context.Context, userID uuid.UUID, name string) (*Secret, error)

	LoadSecretByID(ctx context.Context, ID uuid.UUID) (*Secret, error)

	LoadSecrets(ctx context.Context, userID uuid.UUID) (*[]Secret, error)

	AddTag(ctx context.Context, secretID uuid.UUID, tag string) error

	DeleteTag(ctx context.Context, secretID uuid.UUID, tag string) error

	// Close закрывает соединение с хранилищем.
	Close()
}
