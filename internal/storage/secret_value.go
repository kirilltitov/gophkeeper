package storage

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

// kindsLoadMap contains a mapping of all kinds (see [api.Kinds]) to their respective loading functions.
var kindsLoadMap = map[api.Kind]func(ctx context.Context, tx pgx.Tx, secret *Secret) (SecretValue, error){
	api.KindCredentials: loadSecretCredentialsValue,
	api.KindNote:        loadSecretNoteValue,
	api.KindBlob:        loadSecretBlobValue,
	api.KindBankCard:    loadSecretBankCardValue,
}

// SecretValue is an interface defining all common methods for all kinds of secrets (see [api.Kinds]).
type SecretValue interface {
	SetID(id uuid.UUID)
	CreateValue(ctx context.Context, tx pgx.Tx, secret *Secret) error
	Kind() api.Kind
}

// EditSecretCredentials edits secret credentials with new values.
func (s *PgSQL) EditSecretCredentials(ctx context.Context, secret *Secret, login string, password string) error {
	if secret.Kind != api.KindCredentials {
		return ErrWrongKind
	}

	query := `update public.secret_credentials set login = $1, password = $2 where id = $3`
	_, err := s.Conn.Exec(ctx, query, login, password, secret.ID)
	return err
}

// EditSecretNote edits secret note with new values.
func (s *PgSQL) EditSecretNote(ctx context.Context, secret *Secret, body string) error {
	if secret.Kind != api.KindNote {
		return ErrWrongKind
	}

	query := `update public.secret_note set body = $1 where id = $2`
	_, err := s.Conn.Exec(ctx, query, body, secret.ID)
	return err
}

// EditSecretBlob edits secret blob with new values.
func (s *PgSQL) EditSecretBlob(ctx context.Context, secret *Secret, body string) error {
	if secret.Kind != api.KindBlob {
		return ErrWrongKind
	}

	query := `update public.secret_blob set body = $1 where id = $2`
	_, err := s.Conn.Exec(ctx, query, body, secret.ID)
	return err
}

// EditSecretBankCard edits secret bank card with new values.
func (s *PgSQL) EditSecretBankCard(ctx context.Context, secret *Secret, name, number, date, cvv string) error {
	if secret.Kind != api.KindBankCard {
		return ErrWrongKind
	}

	query := `update public.secret_bank_card set name = $1, number = $2, date = $3, cvv = $4 where id = $5`
	_, err := s.Conn.Exec(ctx, query, name, number, date, cvv, secret.ID)
	return err
}

// CreateValue creates a new secret value.
func (s *SecretBankCard) CreateValue(ctx context.Context, tx pgx.Tx, secret *Secret) error {
	if secret.Kind != api.KindBankCard {
		return ErrWrongKind
	}

	query := `insert into public.secret_bank_card (id, name, number, date, cvv) values ($1, $2, $3, $4, $5)`
	_, err := tx.Exec(ctx, query, s.ID, s.Name, s.Number, s.Date, s.CVV)
	return err
}

// CreateValue creates a new secret value.
func (s *SecretBlob) CreateValue(ctx context.Context, tx pgx.Tx, secret *Secret) error {
	if secret.Kind != api.KindBlob {
		return ErrWrongKind
	}

	query := `insert into public.secret_blob (id, body) values ($1, $2)`
	_, err := tx.Exec(ctx, query, s.ID, s.Body)
	return err
}

// CreateValue creates a new secret value.
func (s *SecretNote) CreateValue(ctx context.Context, tx pgx.Tx, secret *Secret) error {
	if secret.Kind != api.KindNote {
		return ErrWrongKind
	}

	query := `insert into public.secret_note (id, body) values ($1, $2)`
	_, err := tx.Exec(ctx, query, s.ID, s.Body)
	return err
}

// CreateValue creates a new secret value.
func (s *SecretCredentials) CreateValue(ctx context.Context, tx pgx.Tx, secret *Secret) error {
	if secret.Kind != api.KindCredentials {
		return ErrWrongKind
	}

	query := `insert into public.secret_credentials (id, login, password) values ($1, $2, $3)`
	_, err := tx.Exec(ctx, query, s.ID, s.Login, s.Password)
	return err
}

// SetID sets parent secret ID to secret value.
func (s *SecretBankCard) SetID(id uuid.UUID) {
	s.ID = id
}

// SetID sets parent secret ID to secret value.
func (s *SecretBlob) SetID(id uuid.UUID) {
	s.ID = id
}

// SetID sets parent secret ID to secret value.
func (s *SecretNote) SetID(id uuid.UUID) {
	s.ID = id
}

// SetID sets parent secret ID to secret value.
func (s *SecretCredentials) SetID(id uuid.UUID) {
	s.ID = id
}

// Kind returns a kind of current secret value.
func (s *SecretBankCard) Kind() api.Kind {
	return api.KindBankCard
}

// Kind returns a kind of current secret value.
func (s *SecretBlob) Kind() api.Kind {
	return api.KindBlob
}

// Kind returns a kind of current secret value.
func (s *SecretNote) Kind() api.Kind {
	return api.KindNote
}

// Kind returns a kind of current secret value.
func (s *SecretCredentials) Kind() api.Kind {
	return api.KindCredentials
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
