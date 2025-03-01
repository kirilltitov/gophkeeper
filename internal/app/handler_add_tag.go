package app

import (
	"errors"
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/pkg/api"
)

// HandlerAddTag adds a tag to a given secret.
//
// Example request:
//
// POST /api/secret/tag/{ID}
//
//	{
//		"tag": "some tag"
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
func (a *Application) HandlerAddTag(w http.ResponseWriter, r *http.Request) {
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

	var req api.TagRequest

	defer r.Body.Close()
	err = parseRequest(w, r.Body, &req)
	if err != nil {
		return
	}

	err = a.Gophkeeper.AddTag(
		ctx,
		*secretID,
		req.Tag,
	)
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
