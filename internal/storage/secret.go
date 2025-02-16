package storage

import "github.com/google/uuid"

type Kind string

const (
	KindCredentials = "credentials"
	KindNote        = "note"
	KindBlob        = "blob"
	KindBankCard    = "bank_card"
)

type Secret struct {
	ID     uuid.UUID `db:"id"`
	UserID uuid.UUID `db:"user_id"`
	name   string    `db:"name"`
	kind   Kind      `db:"kind"`
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
