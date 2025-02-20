package utils

import (
	"context"

	"github.com/google/uuid"
)

type CtxUserIDKey struct{}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(CtxUserIDKey{}).(uuid.UUID)
	return userID, ok
}

func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, CtxUserIDKey{}, userID)
}
