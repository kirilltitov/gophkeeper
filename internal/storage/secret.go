package storage

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

// Secret is a root entity of the service containing all secret fields.
type Secret struct {
	ID          uuid.UUID   `db:"id" json:"id"`                     // ID is a unique identifier.
	UserID      uuid.UUID   `db:"user_id" json:"user_id"`           // UserID is the secret owner's identifier.
	Name        string      `db:"name" json:"name"`                 // Name is secret name.
	Tags        Tags        `db:"tags" json:"tags"`                 // Tags is a list of secret tags.
	Kind        api.Kind    `db:"kind" json:"kind"`                 // Kind is a kind of secret (see [api.Kinds]).
	IsEncrypted bool        `db:"is_encrypted" json:"is_encrypted"` // IsEncrypted indicates whether secret is encrypted.
	Value       SecretValue `json:"value"`                          // Value is actual secret value (depending on kind).
}

// SecretCredentials is a model containing secret credentials values.
type SecretCredentials struct {
	ID       uuid.UUID `db:"id" json:"id"`             // ID is a unique secret identifier.
	Login    string    `db:"login" json:"login"`       // Login is credentials login.
	Password string    `db:"password" json:"password"` // Password is credentials password.
}

// SecretNote is a model containing secret note value.
type SecretNote struct {
	ID   uuid.UUID `db:"id" json:"id"`     // ID is a unique secret identifier.
	Body string    `db:"body" json:"body"` // Body is note body.
}

// SecretBlob is a model containing secret blob value.
type SecretBlob struct {
	ID   uuid.UUID `db:"id" json:"id"`     // ID is a unique secret identifier.
	Body string    `db:"body" json:"body"` // Body is blob body (in ASCII form, preferably base64).
}

// SecretBankCard is a model containing secret bank card values.
type SecretBankCard struct {
	ID     uuid.UUID `db:"id" json:"id"`         // ID is a unique secret identifier.
	Name   string    `db:"name" json:"name"`     // Name is cardholder name.
	Number string    `db:"number" json:"number"` // Number is card number.
	Date   string    `db:"date" json:"date"`     // Date is card expiration date.
	CVV    string    `db:"cvv" json:"cvv"`       // CVV is CVV (or CVC).
}

// CreateSecret creates a new secret in DB.
func (s *PgSQL) CreateSecret(ctx context.Context, secret *Secret) error {
	return WithVoidTransaction(ctx, s, func(tx pgx.Tx) error {
		_, ok := api.Kinds[secret.Kind]
		if !ok {
			return ErrInvalidKind
		}

		if err := createSecret(ctx, tx, secret); err != nil {
			return err
		}

		return tx.Commit(ctx)
	})
}

// DeleteSecret deletes a secret from a DB.
func (s *PgSQL) DeleteSecret(ctx context.Context, secretID uuid.UUID) error {
	query := `delete from public.secret where id = $1`
	_, err := s.Conn.Exec(ctx, query, secretID)
	return err
}

// LoadSecretByName loads a secret by name.
func (s *PgSQL) LoadSecretByName(ctx context.Context, userID uuid.UUID, name string) (*Secret, error) {
	return WithTransaction(ctx, s, func(tx pgx.Tx) (*Secret, error) {
		result, err := loadSecretByName(ctx, tx, userID, name)

		if err != nil {
			return nil, err
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}

		return result, err
	})
}

// LoadSecretByID loads a secret by ID.
func (s *PgSQL) LoadSecretByID(ctx context.Context, secretID uuid.UUID) (*Secret, error) {
	return WithTransaction(ctx, s, func(tx pgx.Tx) (*Secret, error) {
		result, err := loadSecretByID(ctx, tx, secretID)

		if err != nil {
			return nil, err
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}

		return result, err
	})
}

// LoadSecrets loads all secrets for given user.
func (s *PgSQL) LoadSecrets(ctx context.Context, userID uuid.UUID) (*[]Secret, error) {
	var rows []Secret

	query := `
		select
			s.*,
			json_agg_strict(t.text) tags
		from secret s
		left join tag t on s.id = t.secret_id
		where s.user_id = $1
		group by s.id
		order by s.name
	`
	err := pgxscan.Select(ctx, s.Conn, &rows, query, userID)
	if err != nil {
		return nil, err
	}

	return &rows, nil
}

// RenameSecret renames secret.
func (s *PgSQL) RenameSecret(ctx context.Context, secretID uuid.UUID, name string) error {
	query := `update public.secret set name = $1 where id = $2`
	_, err := s.Conn.Exec(ctx, query, name, secretID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateSecretFound
		}
		return err
	}

	return nil
}

func loadSecretByName(ctx context.Context, tx pgx.Tx, userID uuid.UUID, name string) (*Secret, error) {
	var secret Secret

	query := `
		select
			s.*,
			json_agg_strict(t.text) tags
		from secret s
		left join tag t on s.id = t.secret_id
		where s.user_id = $1 and s.name = $2
		group by s.id
	`
	if err := pgxscan.Get(ctx, tx, &secret, query, userID, name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	secretValue, err := loadSecretValue(ctx, tx, &secret)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Log.Errorf(
				"Missing secret value '%s' for secret %s, this MUST NOT ever happen",
				secret.Kind,
				secret.ID.String(),
			)
		}
		return nil, err
	}
	secret.Value = secretValue

	return &secret, nil
}

func loadSecretByID(ctx context.Context, tx pgx.Tx, secretID uuid.UUID) (*Secret, error) {
	var secret Secret

	query := `
		select
			s.*,
			json_agg_strict(t.text) tags
		from secret s
		left join tag t on s.id = t.secret_id
		where s.id = $1
		group by s.id
	`
	if err := pgxscan.Get(ctx, tx, &secret, query, secretID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	secretValue, err := loadSecretValue(ctx, tx, &secret)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			utils.Log.Errorf(
				"Missing secret value '%s' for secret %s, this MUST NOT ever happen",
				secret.Kind,
				secret.ID.String(),
			)
		}
		return nil, err
	}
	secret.Value = secretValue

	return &secret, nil
}

func createSecret(ctx context.Context, tx pgx.Tx, secret *Secret) error {
	query := `insert into public.secret (id, user_id, name, kind, is_encrypted) values ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(ctx, query, secret.ID, secret.UserID, secret.Name, secret.Kind, secret.IsEncrypted)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateSecretFound
		}
		return err
	}

	if err := secret.Value.CreateValue(ctx, tx, secret); err != nil {
		return err
	}

	return nil
}
