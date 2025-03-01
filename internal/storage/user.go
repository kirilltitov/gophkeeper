package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// User is a user entity containing all user fields.
type User struct {
	ID        uuid.UUID `db:"id"`         // ID is a unique user identifier.
	Login     string    `db:"login"`      // Login is a unique user login.
	Password  string    `db:"password"`   // Password is user's hashed password.
	CreatedAt time.Time `db:"created_at"` // CreatedAt is a date of user creation.
}

// IsValidPassword returns true if given raw password is equal to hashed user password upon hashing.
func (u *User) IsValidPassword(password string) bool {
	return u.getHashedPassword(password) == u.Password
}

func (u *User) getHashedPassword(rawPassword string) string {
	dt := strconv.FormatInt(u.CreatedAt.Unix(), 10)

	h := sha256.New()
	_, _ = io.WriteString(h, rawPassword)
	_, _ = io.WriteString(h, dt)

	result := hex.EncodeToString(h.Sum(nil))

	utils.Log.Debugf("Hashing sha256('%s' + '%s') = '%s'", rawPassword, dt, result)

	return result
}

// NewUser creates and returns a new configured user.
func NewUser(id uuid.UUID, login string, rawPassword string) User {
	user := User{
		ID:        id,
		Login:     login,
		CreatedAt: time.Now(),
	}

	user.Password = user.getHashedPassword(rawPassword)

	return user
}

func loadUser(ctx context.Context, tx pgx.Tx, login string) (*User, error) {
	var row User

	if err := pgxscan.Get(ctx, tx, &row, `select * from public.user where login = $1`, login); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return &row, nil
}

func createUser(ctx context.Context, tx pgx.Tx, user User) error {
	query := `insert into public.user (id, login, password, created_at) values ($1, $2, $3, $4)`
	_, err := tx.Exec(ctx, query, user.ID, user.Login, user.Password, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// LoadUser loads a user from DB for given login.
func (s *PgSQL) LoadUser(ctx context.Context, login string) (*User, error) {
	return WithTransaction(ctx, s, func(tx pgx.Tx) (*User, error) {
		user, err := loadUser(ctx, tx, login)
		if err != nil {
			return nil, err
		}

		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}

		return user, nil
	})
}

// CreateUser creates a new user in DB.
func (s *PgSQL) CreateUser(ctx context.Context, user User) error {
	return WithVoidTransaction(ctx, s, func(tx pgx.Tx) error {
		if err := createUser(ctx, tx, user); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return ErrDuplicateUserFound
			}
			return err
		}

		return tx.Commit(ctx)
	})
}
