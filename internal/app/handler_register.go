package app

import (
	"errors"
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/storage"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

func (a *Application) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := utils.Log

	defer r.Body.Close()
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := parseRequest(w, r.Body, &req); err != nil {
		return
	}

	user, err := a.Gophkeeper.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		log.Errorf("Error while registering new user: %v", err)
		var code int
		switch {
		case errors.Is(err, storage.ErrDuplicateUserFound):
			code = http.StatusConflict
		case errors.Is(err, gophkeeper.ErrEmptyLogin), errors.Is(err, gophkeeper.ErrEmptyPassword):
			code = http.StatusBadRequest
		default:
			code = http.StatusInternalServerError
		}
		returnErrorWithCode(w, code, "could not register")
		return
	}

	cookie, err := a.CreateAuthCookie(*user)
	if err != nil {
		log.Errorf("Error while issuing cookie: %v", err)
		returnErrorWithCode(w, http.StatusInternalServerError, "could not register")
		return
	}

	http.SetCookie(w, cookie)

	returnSuccessWithCode(w, http.StatusOK, nil)
}
