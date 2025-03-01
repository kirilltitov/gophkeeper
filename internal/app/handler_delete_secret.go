package app

import (
	"errors"
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// HandlerDeleteSecret deletes a secret.
//
// Example request:
//
// DELETE /api/secret/{ID}
//
// May response with codes 200, 401, 500.
func (a *Application) HandlerDeleteSecret(w http.ResponseWriter, r *http.Request) {
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

	err = a.Gophkeeper.DeleteSecret(ctx, *secretID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, gophkeeper.ErrNoAuth) {
			code = http.StatusUnauthorized
		}
		returnErrorWithCode(w, code, err.Error())
		return
	}

	returnEmptySuccessWithCode(w, http.StatusOK)
}
