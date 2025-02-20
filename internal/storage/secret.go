package storage

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kirilltitov/gophkeeper/internal/utils"
)

type Kind string

const (
	KindCredentials = "credentials"
	KindNote        = "note"
	KindBlob        = "blob"
	KindBankCard    = "bank_card"
)

var Kinds = map[Kind]bool{
	KindCredentials: true,
	KindNote:        true,
	KindBlob:        true,
	KindBankCard:    true,
}

type Secret struct {
	ID     uuid.UUID `db:"id"`
	UserID uuid.UUID `db:"user_id"`
	Name   string    `db:"name"`
	Tags   Tags      `db:"tags"`
	Kind   Kind      `db:"kind"`

	value SecretValue
}

type SecretCredentials struct {
	ID       uuid.UUID `db:"id"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
}

type SecretNote struct {
	ID   uuid.UUID `db:"id"`
	Body string    `db:"body"`
}

type SecretBlob struct {
	ID   uuid.UUID `db:"id"`
	Body string    `db:"body"`
}

type SecretBankCard struct {
	ID     uuid.UUID `db:"id"`
	Name   string    `db:"name"`
	Number string    `db:"number"`
	Date   string    `db:"date"`
	CVV    string    `db:"cvv"`
}

func (s *Secret) Value() SecretValue {
	return s.value
}

func (s *PgSQL) CreateSecret(ctx context.Context, secret *Secret) error {
	return WithVoidTransaction(ctx, s, func(tx pgx.Tx) error {
		_, ok := Kinds[secret.Kind]
		if !ok {
			return ErrInvalidKind
		}

		if err := createSecret(ctx, tx, secret); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return ErrDuplicateSecretFound
			}
			return err
		}

		if err := tx.Commit(ctx); err != nil {
			return err
		}

		return nil
	})
}

func (s *PgSQL) DeleteSecret(ctx context.Context, secretID uuid.UUID) error {
	query := `delete from public.secret where id = $1`
	_, err := s.Conn.Exec(ctx, query, secretID)
	return err
}

func (s *PgSQL) LoadSecretByName(ctx context.Context, userID uuid.UUID, name string) (*Secret, error) {
	return WithTransaction(ctx, s, func(tx pgx.Tx) (*Secret, error) {
		return loadSecretByName(ctx, tx, userID, name)
	})
}

func (s *PgSQL) LoadSecretByID(ctx context.Context, ID uuid.UUID) (*Secret, error) {
	return WithTransaction(ctx, s, func(tx pgx.Tx) (*Secret, error) {
		return loadSecretByID(ctx, tx, ID)
	})
}

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
	secret.value = secretValue

	return &secret, nil
}

func loadSecretByID(ctx context.Context, tx pgx.Tx, ID uuid.UUID) (*Secret, error) {
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
	if err := pgxscan.Get(ctx, tx, &secret, query, ID); err != nil {
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
	secret.value = secretValue

	return &secret, nil
}

func createSecret(ctx context.Context, tx pgx.Tx, secret *Secret) error {
	query := `insert into public.secret (id, user_id, name, kind) values ($1, $2, $3, $4)`
	_, err := tx.Exec(ctx, query, secret.ID, secret.UserID, secret.Name, secret.Kind)
	if err != nil {
		return err
	}

	if err := secret.value.CreateValue(ctx, tx, secret); err != nil {
		return err
	}

	return nil
}
