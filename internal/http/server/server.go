// Package server - классический http сервер
//
// # Описание
//
// Конфигурация и запуск HTTP сервера.
package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/http/router"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/shutdown"
	"github.com/Alheor/shorturl/internal/tlscerts"

	"go.uber.org/zap"
)

// StartServer запуск http сервера
func StartServer(cfg *config.Options) {

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: router.GetRoutes(),
	}

	shutdown.GetCloser().Add(func(ctx context.Context) error {
		err := srv.Shutdown(ctx)
		if err != nil {
			logger.Error(`error while shutting down server`, err)
		}

		return nil
	})

	go func() {
		if cfg.EnableHTTPS {
			logger.Info("Starting HTTPS server", zap.String("addr", cfg.Addr))

			// Создаем TLS конфигурацию для безопасности
			srv.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				},
			}

			var certFile, keyFile string
			var err error

			if cfg.TLSCert != `` && cfg.TLSKey != `` {
				certFile, keyFile, err = tlscerts.LoadCert(cfg.TLSCert, cfg.TLSKey)
			} else {
				certFile, keyFile, err = tlscerts.GenerateCert()
			}

			if err != nil {
				logger.Fatal(`error while prepare certificates`, err)
			}

			if err = srv.ListenAndServeTLS(certFile, keyFile); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal(`error while starting https server`, err)
			}

		} else {
			logger.Info("Starting HTTP server", zap.String("addr", cfg.Addr))

			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal(`error while starting http server`, err)
			}
		}
	}()
}
