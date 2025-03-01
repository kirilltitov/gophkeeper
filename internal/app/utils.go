package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/kirilltitov/gophkeeper/pkg/api"
)

func parseRequest(w http.ResponseWriter, r io.Reader, target any) error {
	var buf bytes.Buffer

	if n, err := buf.ReadFrom(r); err != nil || n == 0 {
		w.WriteHeader(http.StatusBadRequest)
		returnErrorWithCode(w, http.StatusBadRequest, "no body")
		if err == nil {
			err = errors.New("no body")
		}
		return err
	}
	if err := json.Unmarshal(buf.Bytes(), &target); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		returnErrorWithCode(w, http.StatusBadRequest, "invalid input JSON")
		return err
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(target); err != nil {
		returnErrorWithCode(w, http.StatusBadRequest, "invalid input JSON")
		return err
	}

	return nil
}

func returnErrorWithCode(w http.ResponseWriter, code int, err string) {
	var resultErr *string
	if err == "" {
		resultErr = nil
	} else {
		resultErr = &err
	}

	returnWithCode(
		w,
		code,
		api.BaseResponse[struct{}]{
			Success: false,
			Error:   resultErr,
			Result:  nil,
		},
	)
}

func returnSuccessWithCode[R any](w http.ResponseWriter, code int, body *R) {
	returnWithCode(
		w,
		code,
		api.BaseResponse[R]{
			Success: true,
			Result:  body,
		},
	)
}

func returnEmptySuccessWithCode(w http.ResponseWriter, code int) {
	returnSuccessWithCode[any](w, code, nil)
}

func returnWithCode(w http.ResponseWriter, code int, body any) {
	w.WriteHeader(code)

	if body != nil {
		responseBytes, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(responseBytes)
		if err != nil {
			panic(err)
		}
	}
}

func getUUIDFromRequest(r *http.Request, key string) (*uuid.UUID, error) {
	idString := chi.URLParam(r, key)
	if idString == "" {
		return nil, errors.New("no " + key)
	}

	result, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
