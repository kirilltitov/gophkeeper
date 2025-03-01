package gophkeeper

import (
	"context"

	"github.com/google/uuid"
	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

// EditSecretCredentials edits existing secret credentials.
func (g *Gophkeeper) EditSecretCredentials(
	ctx context.Context,
	secretID uuid.UUID,
	login, password string,
) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	if secret.Kind != api.KindCredentials {
		return storage.ErrWrongKind
	}

	return g.Container.Storage.EditSecretCredentials(ctx, secret, login, password)
}

// EditSecretNote edits existing secret note.
func (g *Gophkeeper) EditSecretNote(
	ctx context.Context,
	secretID uuid.UUID,
	body string,
) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	if secret.Kind != api.KindNote {
		return storage.ErrWrongKind
	}

	return g.Container.Storage.EditSecretNote(ctx, secret, body)
}

// EditSecretBlob edits existing secret blob.
func (g *Gophkeeper) EditSecretBlob(
	ctx context.Context,
	secretID uuid.UUID,
	body string,
) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	if secret.Kind != api.KindBlob {
		return storage.ErrWrongKind
	}

	return g.Container.Storage.EditSecretBlob(ctx, secret, body)
}

// EditSecretBankCard edits existing secret bank card.
func (g *Gophkeeper) EditSecretBankCard(
	ctx context.Context,
	secretID uuid.UUID,
	name, number, date, cvv string,
) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	if secret.Kind != api.KindBankCard {
		return storage.ErrWrongKind
	}

	return g.Container.Storage.EditSecretBankCard(ctx, secret, name, number, date, cvv)
}
