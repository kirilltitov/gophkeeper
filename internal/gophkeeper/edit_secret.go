package gophkeeper

import (
	"context"

	"github.com/google/uuid"
	"github.com/kirilltitov/gophkeeper/internal/storage"
)

func (g *Gophkeeper) EditSecretCredentials(
	ctx context.Context,
	secretID uuid.UUID,
	login, password string,
) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	if secret.Kind != storage.KindCredentials {
		return storage.ErrWrongKind
	}

	return g.Container.Storage.EditSecretCredentials(ctx, secret, login, password)
}

func (g *Gophkeeper) EditSecretNote(
	ctx context.Context,
	secretID uuid.UUID,
	body string,
) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	if secret.Kind != storage.KindNote {
		return storage.ErrWrongKind
	}

	return g.Container.Storage.EditSecretNote(ctx, secret, body)
}

func (g *Gophkeeper) EditSecretBlob(
	ctx context.Context,
	secretID uuid.UUID,
	body string,
) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	if secret.Kind != storage.KindBlob {
		return storage.ErrWrongKind
	}

	return g.Container.Storage.EditSecretBlob(ctx, secret, body)
}

func (g *Gophkeeper) EditSecretBankCard(
	ctx context.Context,
	secretID uuid.UUID,
	name, number, date, cvv string,
) error {
	secret, err := g.loadSecretAndAuthorize(ctx, secretID)
	if err != nil {
		return err
	}

	if secret.Kind != storage.KindBankCard {
		return storage.ErrWrongKind
	}

	return g.Container.Storage.EditSecretBankCard(ctx, secret, name, number, date, cvv)
}
