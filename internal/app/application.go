package app

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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
	router := chi.NewRouter()

	router.Use(utils.WithLogging)

	router.Mount("/debug", middleware.Profiler())

	router.Post("/api/login", a.HandlerLogin)
	router.Post("/api/register", a.HandlerRegister)

	return router
}
