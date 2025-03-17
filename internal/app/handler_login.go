package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// HandlerLogin performs login with given login and password.
//
// Example request:
//
// POST /api/login
//
//	{
//		"login":    "john.appleseed",
//		"password": "MoolyFTW",
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
// Also set a JWT cookie on success.
//
// May response with codes 200, 401, 500.
func (a *Application) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	log := utils.Log

	defer r.Body.Close()

	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := parseRequest(w, r.Body, &req); err != nil {
		return
	}

	fmt.Printf("about to login user '%s' with password '%s'", req.Login, req.Password) // todo remove
	user, err := a.Gophkeeper.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		log.Errorf("Error while logging in: %v", err)
		var code int
		switch {
		case errors.Is(err, gophkeeper.ErrAuthFailed):
			code = http.StatusUnauthorized
		default:
			code = http.StatusInternalServerError
		}
		returnErrorWithCode(w, code, "could not authenticate")
		return
	}

	cookie, err := a.CreateAuthCookie(*user)
	if err != nil {
		log.Errorf("Error while issuing cookie: %v", err)
		returnErrorWithCode(w, http.StatusInternalServerError, "could not authenticate")
		return
	}

	http.SetCookie(w, cookie)

	returnEmptySuccessWithCode(w, http.StatusOK)
}
