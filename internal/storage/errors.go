package storage

import (
	"errors"
)

// ErrNotFound is an error indicating that certain entity is not found in DB.
var ErrNotFound = errors.New("not found")

// ErrDuplicateUserFound is an error indicating that user with given login already exists.
var ErrDuplicateUserFound = errors.New("user with this login already exists")

// ErrDuplicateSecretFound is an error indicating that secret with given name already exists.
var ErrDuplicateSecretFound = errors.New("secret with this name already exists")

// ErrInvalidKind is an error indicating that secret kind is unknown.
var ErrInvalidKind = errors.New("invalid secret kind")

// ErrWrongKind is an error indicating that secret kind and factual value differ.
var ErrWrongKind = errors.New("secret kind does not match actual secret value")
