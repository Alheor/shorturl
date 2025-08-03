// Package interceptor - gRPC интерцепторы
//
// # Описание
//
// Реализует интерцепторы для gRPC сервера, включая аутентификацию пользователей.
package interceptor

import (
	"context"
	"strings"

	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/userauth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor интерцептор для аутентификации пользователей
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Методы, которые не требуют аутентификации
	publicMethods := map[string]bool{
		"/shorturl.v1.ShortURLService/Ping":                              true,
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":      true,
	}

	if publicMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	// Извлекаем метаданные из контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	// Ищем токен аутентификации в заголовках
	var authToken string
	if values := md.Get("authorization"); len(values) > 0 {
		authToken = strings.TrimPrefix(values[0], "Bearer ")

	}

	// Альтернативно ищем в заголовке x-auth-user
	if authToken == "" {
		if values := md.Get("x-auth-user"); len(values) > 0 {
			authToken = values[0]
		}
	}

	if authToken == "" {
		return nil, status.Error(codes.Unauthenticated, "missing auth token")
	}

	// Попытка разпарсить токен как cookie
	user, err := userauth.ParseCookieToken(authToken)
	if err != nil {
		logger.Error("Failed to parse auth token", err)
		return nil, status.Error(codes.Unauthenticated, "invalid auth token")
	}

	// Добавляем пользователя в контекст
	ctx = context.WithValue(ctx, models.ContextValueName, user)

	return handler(ctx, req)
}

// StreamAuthInterceptor интерцептор для аутентификации в стриминговых методах
func StreamAuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// Методы, которые не требуют аутентификации
	publicMethods := map[string]bool{
		"/shorturl.v1.ShortURLService/Ping":                              true,
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": true,
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":      true,
	}

	if publicMethods[info.FullMethod] {
		return handler(srv, ss)
	}

	// Извлекаем метаданные из контекста
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	// Ищем токен аутентификации в заголовках
	var authToken string
	if values := md.Get("authorization"); len(values) > 0 {
		authToken = strings.TrimPrefix(values[0], "Bearer ")
	}

	// Альтернативно ищем в заголовке x-auth-user
	if authToken == "" {
		if values := md.Get("x-auth-user"); len(values) > 0 {
			authToken = values[0]
		}
	}

	if authToken == "" {
		return status.Error(codes.Unauthenticated, "missing auth token")
	}

	// Попытка разпарсить токен как cookie
	user, err := userauth.ParseCookieToken(authToken)
	if err != nil {
		logger.Error("Failed to parse auth token", err)
		return status.Error(codes.Unauthenticated, "invalid auth token")
	}

	// Создаем обертку для ServerStream с обновленным контекстом
	wrappedStream := &serverStreamWithAuth{
		ServerStream: ss,
		ctx:          context.WithValue(ss.Context(), models.ContextValueName, user),
	}

	return handler(srv, wrappedStream)
}

// serverStreamWithAuth обертка для ServerStream с аутентификацией
type serverStreamWithAuth struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStreamWithAuth) Context() context.Context {
	return s.ctx
}
