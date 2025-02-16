package gophkeeper

import (
	"context"

	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// Register Создает нового пользователя с заданным логином и паролем
func (g Gophkeeper) Register(ctx context.Context, login string, rawPassword string) (*storage.User, error) {
	if login == "" {
		return nil, ErrEmptyLogin
	}
	if rawPassword == "" {
		return nil, ErrEmptyPassword
	}

	userID := utils.NewUUID6()
	user := storage.NewUser(userID, login, rawPassword)

	if err := g.Container.Storage.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &user, nil
}
