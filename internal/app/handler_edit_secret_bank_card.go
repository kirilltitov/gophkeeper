package app

import (
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func (a *Application) HandlerEditSecretBankCard(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

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

	var req requestBankCard

	defer r.Body.Close()
	err = parseRequest(w, r.Body, &req)
	if err != nil {
		return
	}

	err = a.Gophkeeper.EditSecretBankCard(
		ctx,
		*secretID,
		req.Name,
		req.Number,
		req.Date,
		req.CVV,
	)
	if err != nil {
		code := http.StatusInternalServerError
		returnErrorWithCode(w, code, err.Error())
		return
	}

	returnSuccessWithCode(w, http.StatusOK, nil)
}
