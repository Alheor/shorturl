// Package server - gRPC сервер
//
// # Описание
//
// Конфигурация и запуск gRPC сервера с поддержкой аналогичных HTTP конфигураций.
package server

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/Alheor/shorturl/internal/config"
	grpcHandler "github.com/Alheor/shorturl/internal/grpc/handler"
	"github.com/Alheor/shorturl/internal/grpc/interceptor"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/shutdown"
	"github.com/Alheor/shorturl/internal/tlscerts"
	pb "github.com/Alheor/shorturl/pkg/shorturl/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// StartGRPCServer запуск gRPC сервера
func StartGRPCServer(cfg *config.Options) {
	if !cfg.EnableGRPC {
		return
	}

	// Создаем listener
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Fatal("failed to listen on gRPC address", err)
	}

	// Создаем gRPC сервер с опциями
	var opts []grpc.ServerOption

	// Настройка TLS для gRPC сервера если включен HTTPS
	if cfg.EnableHTTPS {
		var certFile, keyFile string

		if cfg.TLSCert != "" && cfg.TLSKey != "" {
			certFile, keyFile, err = tlscerts.LoadCert(cfg.TLSCert, cfg.TLSKey)
		} else {
			certFile, keyFile, err = tlscerts.GenerateCert()
		}

		if err != nil {
			logger.Fatal("error while prepare gRPC certificates", err)
		}

		// Загружаем сертификат
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			logger.Fatal("failed to load gRPC certificate", err)
		}

		// Создаем TLS конфигурацию для безопасности
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			},
		}

		// Добавляем TLS credentials
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.Creds(creds))

		logger.Info("gRPC server will use TLS")
	}

	// Добавляем интерцепторы
	opts = append(opts,
		grpc.UnaryInterceptor(interceptor.AuthInterceptor),
		grpc.StreamInterceptor(interceptor.StreamAuthInterceptor),
	)

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer(opts...)

	// Регистрируем наш сервис
	shortURLServer := grpcHandler.NewShortURLServer()
	pb.RegisterShortURLServiceServer(grpcServer, shortURLServer)
	reflection.Register(grpcServer)

	// Добавляем graceful shutdown
	shutdown.GetCloser().Add(func(ctx context.Context) error {
		logger.Info("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		return nil
	})

	// Запускаем сервер в отдельной горутине
	go func() {
		if cfg.EnableHTTPS {
			logger.Info("Starting gRPC server with TLS", zap.String("addr", cfg.GRPCAddr))
		} else {
			logger.Info("Starting gRPC server", zap.String("addr", cfg.GRPCAddr))
		}

		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("gRPC server failed to serve", err)
		}
	}()
}
