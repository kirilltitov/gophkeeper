package utils

import "github.com/google/uuid"

// NewUUID6 is a helper returning a new UUID (v6) instance.
//
// Warning: this method might raise panic should the generation fail (which is an unlikely event),
// so it SHOULD only be used in tests and other non-production code.
func NewUUID6() uuid.UUID {
	result, err := uuid.NewV6()
	if err != nil {
		panic(err)
	}
	return result
}
