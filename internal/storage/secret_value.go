package storage

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"github.com/kirilltitov/gophkeeper/internal/utils"
)

var kindsLoadMap = map[Kind]func(ctx context.Context, tx pgx.Tx, secret *Secret) (SecretValue, error){
	KindCredentials: loadSecretCredentialsValue,
	KindNote:        loadSecretNoteValue,
	KindBlob:        loadSecretBlobValue,
	KindBankCard:    loadSecretBankCardValue,
}

type SecretValue interface {
	CreateValue(ctx context.Context, tx pgx.Tx, secret *Secret) error
}

func (s *SecretBankCard) CreateValue(ctx context.Context, tx pgx.Tx, secret *Secret) error {
	if secret.Kind != KindBankCard {
		return ErrWrongKind
	}

	query := `insert into public.secret_bank_card (id, name, number, date, cvv) values ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(ctx, query, s.ID, s.Name, s.Number, s.Date, s.CVV)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretBlob) CreateValue(ctx context.Context, tx pgx.Tx, secret *Secret) error {
	if secret.Kind != KindBlob {
		return ErrWrongKind
	}

	query := `insert into public.secret_blob (id, body) values ($1, $2)`
	_, err := tx.Exec(ctx, query, s.ID, s.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretNote) CreateValue(ctx context.Context, tx pgx.Tx, secret *Secret) error {
	if secret.Kind != KindNote {
		return ErrWrongKind
	}

	query := `insert into public.secret_note (id, body) values ($1, $2)`
	_, err := tx.Exec(ctx, query, s.ID, s.Body)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretCredentials) CreateValue(ctx context.Context, tx pgx.Tx, secret *Secret) error {
	if secret.Kind != KindCredentials {
		return ErrWrongKind
	}

	query := `insert into public.secret_credentials (id, login, password) values ($1, $2, $3)`
	_, err := tx.Exec(ctx, query, s.ID, s.Login, s.Password)
	if err != nil {
		return err
	}

	return nil
}

func loadSecretValue(ctx context.Context, tx pgx.Tx, secret *Secret) (SecretValue, error) {
	loadFunc, ok := kindsLoadMap[secret.Kind]
	if !ok {
		utils.Log.Errorf("Invalid secret kind '%s' for secret %s", secret.Kind, secret.ID.String())
		return nil, ErrInvalidKind
	}

	return loadFunc(ctx, tx, secret)
}

func loadSecretBankCardValue(ctx context.Context, tx pgx.Tx, secret *Secret) (SecretValue, error) {
	var result SecretBankCard

	err := pgxscan.Get(
		ctx,
		tx,
		&result,
		`select * from public.secret_bank_card where id = $1`,
		secret.ID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return &result, nil
}

func loadSecretBlobValue(ctx context.Context, tx pgx.Tx, secret *Secret) (SecretValue, error) {
	var result SecretBlob

	err := pgxscan.Get(
		ctx,
		tx,
		&result,
		`select * from public.secret_blob where id = $1`,
		secret.ID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return &result, nil
}

func loadSecretNoteValue(ctx context.Context, tx pgx.Tx, secret *Secret) (SecretValue, error) {
	var result SecretNote

	err := pgxscan.Get(
		ctx,
		tx,
		&result,
		`select * from public.secret_note where id = $1`,
		secret.ID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return &result, nil
}

func loadSecretCredentialsValue(ctx context.Context, tx pgx.Tx, secret *Secret) (SecretValue, error) {
	var result SecretCredentials

	err := pgxscan.Get(
		ctx,
		tx,
		&result,
		`select * from public.secret_credentials where id = $1`,
		secret.ID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return &result, nil
}
