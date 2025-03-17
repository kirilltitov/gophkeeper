package app

import (
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

// HandlerEditSecretNote edits a secret note.
//
// Example request:
//
// POST /api/secret/edit/note/{ID}
//
//	{
//		"body": "updated note"
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
func (a *Application) HandlerEditSecretNote(w http.ResponseWriter, r *http.Request) {
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

	var req api.SecretNote

	defer r.Body.Close()
	err = parseRequest(w, r.Body, &req)
	if err != nil {
		return
	}

	err = a.Gophkeeper.EditSecretNote(
		ctx,
		*secretID,
		req.Body,
	)
	if err != nil {
		code := http.StatusInternalServerError
		returnErrorWithCode(w, code, err.Error())
		return
	}

	returnEmptySuccessWithCode(w, http.StatusOK)
}
