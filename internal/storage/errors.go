package storage

import (
	"errors"
)

var ErrNotFound = errors.New("not found")
var ErrDuplicateFound = errors.New("user with this login already exists")
