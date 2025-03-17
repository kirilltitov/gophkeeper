package app

import (
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

// HandlerEditSecretCredentials edits a secret credentials.
//
// Example request:
//
// POST /api/secret/edit/credentials/{ID}
//
//	{
//		"url":      "https://steamcommunity.com/login/home",
//		"login":    "john.appleseed",
//		"password": "MoolyFTW"
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
func (a *Application) HandlerEditSecretCredentials(w http.ResponseWriter, r *http.Request) {
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

	var req api.SecretCredentials

	defer r.Body.Close()
	err = parseRequest(w, r.Body, &req)
	if err != nil {
		return
	}

	err = a.Gophkeeper.EditSecretCredentials(
		ctx,
		*secretID,
		req.URL,
		req.Login,
		req.Password,
	)
	if err != nil {
		code := http.StatusInternalServerError
		returnErrorWithCode(w, code, err.Error())
		return
	}

	returnEmptySuccessWithCode(w, http.StatusOK)
}
