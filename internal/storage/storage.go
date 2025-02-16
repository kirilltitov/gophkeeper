package storage

import "context"

//go:generate mockery
type Storage interface {
	// CreateUser создает нового пользователя в хранилище.
	CreateUser(ctx context.Context, user User) error

	// LoadUser пытается найти пользователя в хранилище по заданному логину.
	LoadUser(ctx context.Context, login string) (*User, error)

	// Close закрывает соединение с хранилищем.
	Close()
}
