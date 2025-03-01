package app

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

// HandlerCreateSecretBlob creates a new secret blob.
// Binary data MUST BE in ASCII form, base64 is preferred.
//
// Example request:
//
// POST /api/secret/create/blob
//
//	{
//		"name": "secret name",
//		"is_encrypted": true,
//		"value": {
//			"body": "0JrQsNC60L7QuS3RgtC+INCx0LXQudC3NjQg0L3QsNC/0YDQuNC80LXRgA=="
//		}
//	}
//
// Example response:
//
//	{
//		"success": true,
//		"result":  {
//	     	"id": "1ee1416c-d537-6ae0-b6c7-0f48c8929427"
//		},
//		"error":   null
//	}
//
// May response with codes 201, 401, 409, 500.
func (a *Application) HandlerCreateSecretBlob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	_, ok := utils.GetUserID(ctx)
	if !ok {
		returnErrorWithCode(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req api.BaseCreateSecretRequest[api.SecretBlob]
	type response struct {
		ID uuid.UUID `json:"id"`
	}

	defer r.Body.Close()
	err = parseRequest(w, r.Body, &req)
	if err != nil {
		return
	}

	secret := &storage.Secret{
		Name:        req.Name,
		IsEncrypted: req.IsEncrypted,
		Value: &storage.SecretBlob{
			Body: req.Value.Body,
		},
	}

	err = a.Gophkeeper.CreateSecret(ctx, secret)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, storage.ErrDuplicateSecretFound) {
			code = http.StatusConflict
		}
		returnErrorWithCode(w, code, err.Error())
		return
	}

	returnSuccessWithCode[api.CreatedSecretResponse](w, http.StatusCreated, &api.CreatedSecretResponse{ID: secret.ID})
}
