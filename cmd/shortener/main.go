package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/httphandler"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/router"
	"github.com/Alheor/shorturl/internal/service"
	"github.com/Alheor/shorturl/internal/shutdown"

	"go.uber.org/zap"
)

var shutdownTimeout = 5 * time.Second

func main() {

	defer func() {
		if err := recover(); err != nil {
			logger.Error(``, err.(error))
			logger.Sync()
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	shutdown.Init()
	cfg := config.Load()

	var err error

	err = logger.Init(nil)
	if err != nil {
		panic(err)
	}

	httphandler.Init(&cfg)
	service.Init(&cfg)

	shutdown.GetCloser().Add(func(ctx context.Context) error {
		logger.Sync()
		return nil
	})

	err = repository.Init(ctx, &cfg, nil)
	if err != nil {
		logger.Fatal(`error while initialize repository`, err)
	}

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: router.GetRoutes(),
	}

	shutdown.GetCloser().Add(srv.Shutdown)

	go func() {
		logger.Info("Starting server", zap.String("addr", cfg.Addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(`error while starting http server`, err)
		}
	}()

	<-ctx.Done()

	println("shutting down ...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	shutdown.GetCloser().Close(shutdownCtx)
}
