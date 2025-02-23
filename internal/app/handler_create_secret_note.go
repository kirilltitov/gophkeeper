package app

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func (a *Application) HandlerCreateSecretNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	_, ok := utils.GetUserID(ctx)
	if !ok {
		returnErrorWithCode(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		Name  string      `json:"name" validate:"required"`
		Value requestNote `json:"value" validate:"required"`
	}
	type response struct {
		ID uuid.UUID `json:"id"`
	}

	defer r.Body.Close()
	err = parseRequest(w, r.Body, &req)
	if err != nil {
		return
	}

	secret := &storage.Secret{
		Name: req.Name,
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

	returnSuccessWithCode(w, http.StatusCreated, response{ID: secret.ID})
}
