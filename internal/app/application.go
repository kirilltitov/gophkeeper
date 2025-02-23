package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// Application является объектом веб-приложения сервиса.
type Application struct {
	Gophkeeper *gophkeeper.Gophkeeper
	Server     *http.Server

	wg *sync.WaitGroup
}

// New создает и возвращает сконфигурированный объект веб-приложения сервиса.
func New(s *gophkeeper.Gophkeeper, wg *sync.WaitGroup) *Application {
	a := &Application{
		Gophkeeper: s,
		Server: &http.Server{
			Addr:              s.Config.ServerAddress,
			ReadHeaderTimeout: time.Second * 5,
		},
		wg: wg,
	}

	a.Server.Handler = utils.GzipHandle(a.createRouter())

	return a
}

// Run запускает веб-сервер приложения.
func (a *Application) Run() {
	defer a.wg.Done()

	logger := utils.Log

	var runFunc func() error

	if a.Gophkeeper.Config.IsTLSEnabled() {
		runFunc = func() error {
			logger.Infof("Starting a HTTPS server at %s", a.Gophkeeper.Config.ServerAddress)
			return a.Server.ListenAndServeTLS(
				"localhost.crt",
				"localhost.key",
			)
		}
	} else {
		runFunc = func() error {
			logger.Warning("Running Gophkeeper without in HTTP mode (without TLS) is extremely unsecure")
			logger.Infof("Starting a HTTP server at %s", a.Gophkeeper.Config.ServerAddress)
			return a.Server.ListenAndServe()
		}
	}

	if err := runFunc(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("HTTP server shutdown")
		} else {
			panic(err)
		}
	}
}

func (a *Application) createRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(utils.WithLogging)

	r.Mount("/debug", middleware.Profiler())

	r.Route("/api", func(r chi.Router) {
		r.Post("/login", a.HandlerLogin)
		r.Post("/register", a.HandlerRegister)

		r.Route("/secret", func(r chi.Router) {
			r.Use(a.WithAuthorization)

			r.Get("/{ID}", a.HandlerGetSecret)
			r.Delete("/{ID}", a.HandlerDeleteSecret)
			r.Post("/{ID}/rename", a.HandlerRenameSecret)

			r.Get("/list", a.HandlerGetSecrets)

			r.Route("/create", func(r chi.Router) {
				r.Post("/bank_card", a.HandlerCreateSecretBankCard)
				r.Post("/credentials", a.HandlerCreateSecretCredentials)
				r.Post("/note", a.HandlerCreateSecretNote)
				r.Post("/blob", a.HandlerCreateSecretBlob)
			})

			r.Route("/edit", func(r chi.Router) {
				r.Post("/bank_card/{ID}", a.HandlerEditSecretBankCard)
				r.Post("/credentials/{ID}", a.HandlerEditSecretCredentials)
				r.Post("/note/{ID}", a.HandlerEditSecretNote)
				r.Post("/blob/{ID}", a.HandlerEditSecretBlob)
			})

			r.Post("/tag/{ID}", a.HandlerAddTag)
			r.Delete("/tag/{ID}", a.HandlerDeleteTag)
		})
	})

	return r
}

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
		baseResponse{
			Success: false,
			Error:   resultErr,
			Result:  nil,
		},
	)
}

func returnSuccessWithCode(w http.ResponseWriter, code int, body any) {
	returnWithCode(
		w,
		code,
		baseResponse{
			Success: true,
			Result:  body,
		},
	)
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
