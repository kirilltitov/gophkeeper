package app

import (
	"errors"
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// HandlerChangeSecretDescription changes an existing secret's description.
//
// Example request:
//
// POST /api/secret/{ID}/change_description
//
//	{
//		"description": "new description"
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
func (a *Application) HandlerChangeSecretDescription(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Description string `json:"description" validate:"required"`
	}

	defer r.Body.Close()
	err = parseRequest(w, r.Body, &req)
	if err != nil {
		return
	}

	err = a.Gophkeeper.ChangeSecretDescription(ctx, *secretID, req.Description)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, storage.ErrDuplicateSecretFound) {
			code = http.StatusConflict
		}
		returnErrorWithCode(w, code, err.Error())
		return
	}

	returnEmptySuccessWithCode(w, http.StatusOK)
}
