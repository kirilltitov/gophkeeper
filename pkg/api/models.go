package api

import "github.com/google/uuid"

// BaseCreateSecretRequest is an envelope for detailed secret response containing base fields and secret Value.
type BaseCreateSecretRequest[V any] struct {
	Name        string `json:"name" validate:"required"`  // Name is secret name.
	IsEncrypted bool   `json:"is_encrypted"`              // IsEncrypted is true if secret value is E2E-encrypted.
	Value       V      `json:"value" validate:"required"` // Value is actual secret value (see [Kinds]).
}

// SecretBankCard is a model representing secret bank card.
type SecretBankCard struct {
	Name   string `json:"name" validate:"required"`   // Name is cardholder name.
	Number string `json:"number" validate:"required"` // Number is card number.
	Date   string `json:"date" validate:"required"`   // Date is card expiration date.
	CVV    string `json:"cvv" validate:"required"`    // CVV is CVV (or CVC).
}

// SecretCredentials is a model representing secret credentials.
type SecretCredentials struct {
	Login    string `json:"login" validate:"required"`    // Login is credentials login.
	Password string `json:"password" validate:"required"` // Password is credentials password.
}

// SecretNote is a model representing secret note.
type SecretNote struct {
	Body string `json:"body" validate:"required"` // Body is note body.
}

// SecretBlob is a model representing secret blob.
type SecretBlob struct {
	Body string `json:"body" validate:"required"` // Body is blob body.
}

// TagRequest is a model representing individual secret tag.
type TagRequest struct {
	Tag string `json:"tag" validate:"required"` // Tag is tag name.
}

// CreatedSecretResponse is a model representing a created secret response.
type CreatedSecretResponse struct {
	ID uuid.UUID `json:"id"` // ID is a unique secret identifier.
}

// Kind is a kind of secret value (see [Kinds]).
type Kind string

const (
	KindCredentials = "credentials" // KindCredentials is representing secret credentials.
	KindNote        = "note"        // KindNote is representing secret note.
	KindBlob        = "blob"        // KindBlob is representing secret blob.
	KindBankCard    = "bank_card"   // KindBankCard is representing secret bank card.
)

// Kinds is a list of all possible kinds of secrets.
var Kinds = map[Kind]bool{
	KindCredentials: true,
	KindNote:        true,
	KindBlob:        true,
	KindBankCard:    true,
}
