package app

import (
	"errors"
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func (a *Application) HandlerGetSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := utils.GetUserID(ctx)
	if !ok {
		returnErrorWithCode(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	secretID, err := getUUIDFromRequest(r, "ID")
	if err != nil {
		returnErrorWithCode(w, http.StatusBadRequest, err.Error())
		return
	}

	secret, err := a.Gophkeeper.GetSecretWithValueByID(ctx, *secretID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, gophkeeper.ErrNoAuth) {
			code = http.StatusUnauthorized
		}
		returnErrorWithCode(w, code, err.Error())
		return
	}

	returnSuccessWithCode(w, http.StatusOK, secret)
}
