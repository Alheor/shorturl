// Сервис сокращения URL адресов.
//
// # Функции сервиса
//
// • сервис принимает длинный URL и возвращает короткий;
//
// • сервис хранит у себя все URL и по запросу возвращает существующий;
//
// • поддерживает работу с пользователями;
//
// • поддерживает массовую загрузку URL для сокращения;
//
// • позволяет получить сразу все сохраненные URL;
//
// • позволяет удалить сохраненный URL.
//
// # Описание сервиса
//
// Для хранения данных, сервис поддерживает работу с базой данных (Postgresql), может хранить данные в файле, а так же в памяти.
//
// Сервис лишен возможности регистрации пользователя.
// В момент обращения он ожидает специальным образом подписанную cookie, по которой попытается авторизовать пользователя.
// Если авторизация не произойдет, то сервис выдаст в ответе новую cookie.
//
// Сервис поддерживает сжатие (Gzip) при взаимодействии по протоколу HTTPS
package main

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/Alheor/shorturl/internal/userauth"

	"go.uber.org/zap"
)

var (
	buildVersion = `N/A`
	buildDate    = `N/A`
	buildCommit  = `N/A`
)

var shutdownTimeout = 5 * time.Second

func main() {

	fmt.Printf("Build version: %s \n", buildVersion)
	fmt.Printf("Build date: %s \n", buildDate)
	fmt.Printf("Build commit: %s \n", buildCommit)

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

	if cfg.SignatureKey == config.DefaultLSignatureKey {
		logger.Error(`Used default signature key! Please change the key!`, nil)
	}

	if len(cfg.SignatureKey) == 0 {
		logger.Fatal(`Signature key is empty`, nil)
	}

	userauth.Init(&cfg)
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
