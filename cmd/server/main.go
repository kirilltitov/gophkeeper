package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/kirilltitov/gophkeeper/internal/app"
	"github.com/kirilltitov/gophkeeper/internal/config"
	"github.com/kirilltitov/gophkeeper/internal/container"
	"github.com/kirilltitov/gophkeeper/internal/gophkeeper"
	"github.com/kirilltitov/gophkeeper/internal/utils"
	"github.com/kirilltitov/gophkeeper/pkg/version"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	v := version.Version{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}
	v.Print(os.Stdout)

	cfg := config.New()
	ctx := context.Background()
	cnt, err := container.New(ctx, cfg)
	if err != nil {
		panic(err)
	}

	service := gophkeeper.New(cfg, cnt)

	run(service)
}

func run(service *gophkeeper.Gophkeeper) {
	wg := &sync.WaitGroup{}
	logger := utils.Log

	application := app.New(service, wg)

	wg.Add(1)
	go application.Run()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-signalChan
	logger.Infof("Received signal: %v", sig)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Shutting down HTTP server")
		if err := application.Server.Shutdown(shutdownCtx); err != nil {
			logger.WithError(err).Error("Could not shutdown HTTP server properly")
		}
		logger.Info("HTTP server is down")
	}()

	wg.Wait()

	service.Container.Storage.Close()

	logger.Info("Goodbye")
}
