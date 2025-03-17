package utils

import (
	"context"

	"github.com/google/uuid"
)

// CtxUserIDKey is a key for setting user ID into [context.Context].
type CtxUserIDKey struct{}

// GetUserID retrieves user ID from given Context.
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(CtxUserIDKey{}).(uuid.UUID)
	return userID, ok
}

// SetUserID sets a user ID to a given Context.
func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, CtxUserIDKey{}, userID)
}
