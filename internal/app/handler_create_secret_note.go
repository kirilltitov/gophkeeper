package app

import (
	"errors"
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

// HandlerCreateSecretNote creates a new secret note (plain text).
//
// Example request:
//
// POST /api/secret/create/note
//
//	{
//		"name": "secret name",
//		"description": "secret description",
//		"is_encrypted": true,
//		"value": {
//			"body": "some secret note"
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
func (a *Application) HandlerCreateSecretNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	_, ok := utils.GetUserID(ctx)
	if !ok {
		returnErrorWithCode(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req api.BaseCreateSecretRequest[api.SecretNote]

	defer r.Body.Close()
	err = parseRequest(w, r.Body, &req)
	if err != nil {
		return
	}

	secret := &storage.Secret{
		Name:        req.Name,
		Description: req.Description,
		IsEncrypted: req.IsEncrypted,
		Value: &storage.SecretNote{
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
