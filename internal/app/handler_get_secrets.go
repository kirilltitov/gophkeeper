package app

import (
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func (a *Application) HandlerGetSecrets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := utils.GetUserID(ctx)
	if !ok {
		returnErrorWithCode(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	secrets, err := a.Gophkeeper.GetSecrets(ctx)
	if err != nil {
		code := http.StatusInternalServerError
		returnErrorWithCode(w, code, err.Error())
		return
	}

	returnSuccessWithCode(w, http.StatusOK, secrets)
}
