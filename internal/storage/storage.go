package storage

import (
	"context"

	"github.com/google/uuid"
)

// Storage is a storage for all entities of the service.
//
//go:generate mockery
type Storage interface {
	// CreateUser creates a new user in DB.
	CreateUser(ctx context.Context, user User) error

	// LoadUser loads a user from DB for given login.
	LoadUser(ctx context.Context, login string) (*User, error)

	// CreateSecret creates a new secret in DB.
	CreateSecret(ctx context.Context, secret *Secret) error

	// RenameSecret renames secret.
	RenameSecret(ctx context.Context, secretID uuid.UUID, name string) error

	// DeleteSecret deletes a secret from a DB.
	DeleteSecret(ctx context.Context, secretID uuid.UUID) error

	// EditSecretCredentials edits secret credentials with new values.
	EditSecretCredentials(ctx context.Context, secret *Secret, login string, password string) error

	// EditSecretNote edits secret note with new values.
	EditSecretNote(ctx context.Context, secret *Secret, body string) error

	// EditSecretBlob edits secret blob with new values.
	EditSecretBlob(ctx context.Context, secret *Secret, body string) error

	// EditSecretBankCard edits secret bank card with new values.
	EditSecretBankCard(ctx context.Context, secret *Secret, name, number, date, cvv string) error

	// LoadSecretByName loads a secret by name.
	LoadSecretByName(ctx context.Context, userID uuid.UUID, name string) (*Secret, error)

	// LoadSecretByID loads a secret by ID.
	LoadSecretByID(ctx context.Context, ID uuid.UUID) (*Secret, error)

	// LoadSecrets loads all secrets for given user.
	LoadSecrets(ctx context.Context, userID uuid.UUID) (*[]Secret, error)

	// AddTag adds a tag to given secret.
	AddTag(ctx context.Context, secretID uuid.UUID, tag string) error

	// DeleteTag removes a tag from given secret.
	DeleteTag(ctx context.Context, secretID uuid.UUID, tag string) error

	// Close закрывает соединение с хранилищем.
	Close()
}
