package app

import (
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

// HandlerEditSecretBankCard edits a secret bank card.
//
// Example request:
//
// POST /api/secret/edit/bank_card/{ID}
//
//	{
//		"name":   "NAME SURNAME",
//		"number": "1234 5678 9012 3456",
//		"date":   "12/34",
//		"cvv":    "322"
//	}
//
// Example response:
//
//	{
//		"success": true,
//		"result":  null,
//		"error":   null
//	}
//
// May response with codes 200, 401, 500.
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

	var req api.SecretBankCard

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

	returnEmptySuccessWithCode(w, http.StatusOK)
}
