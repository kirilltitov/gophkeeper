package app

import (
	"errors"
	"net/http"

	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/utils"
	httpUtils "github.com/kirilltitov/gophkeeper/internal/utils/http"
)

func (a *Application) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log := utils.Log

	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := httpUtils.ParseRequest(w, r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
		w.WriteHeader(code)
		return
	}

	cookie, err := a.CreateAuthCookie(*user)
	if err != nil {
		log.Errorf("Error while issuing cookie: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}
