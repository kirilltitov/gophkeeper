package app

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/utils"
)

// Application is an object holding all application guts.
type Application struct {
	Gophkeeper *gophkeeper.Gophkeeper
	Server     *http.Server

	wg *sync.WaitGroup
}

// New creates and returns a configured instance of [Application].
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

// Run runs application server.
func (a *Application) Run() {
	defer a.wg.Done()

	logger := utils.Log

	var runFunc func() error

	if a.Gophkeeper.Config.IsTLSEnabled() {
		runFunc = func() error {
			logger.Infof("Starting a HTTPS server at %s", a.Gophkeeper.Config.ServerAddress)
			return a.Server.ListenAndServeTLS(
				a.Gophkeeper.Config.TLSCertFile,
				a.Gophkeeper.Config.TLSKeyFile,
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

	r.Route("/api", func(r chi.Router) {
		r.Post("/login", a.HandlerLogin)
		r.Post("/register", a.HandlerRegister)

		r.Route("/secret", func(r chi.Router) {
			r.Use(a.WithAuthorization)

			r.Get("/{ID}", a.HandlerGetSecret)
			r.Delete("/{ID}", a.HandlerDeleteSecret)
			r.Post("/{ID}/rename", a.HandlerRenameSecret)
			r.Post("/{ID}/change_description", a.HandlerChangeSecretDescription)

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
