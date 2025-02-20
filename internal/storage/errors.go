package storage

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrDuplicateUserFound = errors.New("user with this login already exists")
var ErrDuplicateSecretFound = errors.New("secret with this name already exists")
var ErrInvalidKind = errors.New("invalid secret kind")
var ErrWrongKind = errors.New("secret kind does not match actual secret value")
