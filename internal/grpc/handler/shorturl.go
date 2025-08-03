// Package handler - gRPC обработчики для сервиса сокращения URL.
//
// # Описание
//
// Реализует gRPC интерфейс как фасад к общему сервисному слою.
// Все обработчики функционируют идентично HTTP аналогам.
package handler

import (
	"context"
	"errors"
	"net/url"

	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/service"
	"github.com/Alheor/shorturl/internal/userauth"
	pb "github.com/Alheor/shorturl/pkg/shorturl/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ShortURLServer реализует gRPC интерфейс ShortURLServiceServer
type ShortURLServer struct {
	pb.UnimplementedShortURLServiceServer
}

// NewShortURLServer создает новый экземпляр gRPC сервера
func NewShortURLServer() *ShortURLServer {
	return &ShortURLServer{}
}

// AddShorten добавляет один URL и возвращает его сокращенную версию
func (s *ShortURLServer) AddShorten(ctx context.Context, req *pb.AddShortenRequest) (*pb.AddShortenResponse, error) {
	logger.Info(`Used gRPC "AddShorten" handler`)

	if req.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "url required")
	}

	if _, err := url.ParseRequestURI(req.Url); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid URL")
	}

	user := userauth.GetUser(ctx)
	if user == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}
	shortURL, err := service.Add(ctx, user, req.Url)
	if err != nil {
		var uniqErr *models.UniqueErr
		if errors.As(err, &uniqErr) {
			return &pb.AddShortenResponse{
				ShortUrl: uniqErr.ShortKey,
			}, status.Error(codes.AlreadyExists, "URL already exists")
		}

		logger.Error("gRPC AddShorten error: ", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.AddShortenResponse{
		ShortUrl: shortURL,
	}, nil
}

// GetURL получает оригинальный URL по сокращенной версии
func (s *ShortURLServer) GetURL(ctx context.Context, req *pb.GetURLRequest) (*pb.GetURLResponse, error) {
	logger.Info(`Used gRPC "GetURL" handler`)

	if req.ShortName == "" {
		return nil, status.Error(codes.InvalidArgument, "short_name required")
	}

	user := userauth.GetUser(ctx)
	if user == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}
	originalURL, isRemoved := service.Get(ctx, user, req.ShortName)
	if originalURL == "" {
		return nil, status.Error(codes.NotFound, "URL not found")
	}

	return &pb.GetURLResponse{
		Url:       originalURL,
		IsRemoved: isRemoved,
	}, nil
}

// AddShortenBatch добавляет несколько URL и возвращает их сокращенные версии
func (s *ShortURLServer) AddShortenBatch(ctx context.Context, req *pb.AddShortenBatchRequest) (*pb.AddShortenBatchResponse, error) {
	logger.Info(`Used gRPC "AddShortenBatch" handler`)

	if len(req.Urls) == 0 {
		return nil, status.Error(codes.InvalidArgument, "urls required")
	}

	user := userauth.GetUser(ctx)
	if user == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	// Преобразуем protobuf структуры в внутренние модели
	batch := make([]models.APIBatchRequestEl, 0, len(req.Urls))
	for _, el := range req.Urls {
		if el.OriginalUrl == "" {
			return nil, status.Error(codes.InvalidArgument, "original_url required")
		}
		if _, err := url.ParseRequestURI(el.OriginalUrl); err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid URL: "+el.OriginalUrl)
		}

		batch = append(batch, models.APIBatchRequestEl{
			CorrelationID: el.CorrelationId,
			OriginalURL:   el.OriginalUrl,
		})
	}

	// user уже получен из контекста выше
	response, err := service.AddBatch(ctx, user, batch)
	if err != nil {
		logger.Error("gRPC AddShortenBatch error: ", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	// Преобразуем ответ в protobuf структуры
	urls := make([]*pb.BatchResponseElement, 0, len(response))
	for _, el := range response {
		urls = append(urls, &pb.BatchResponseElement{
			CorrelationId: el.CorrelationID,
			ShortUrl:      el.ShortURL,
		})
	}

	return &pb.AddShortenBatchResponse{
		Urls: urls,
	}, nil
}

// GetUserURLs получает все URL пользователя (стриминговый ответ)
func (s *ShortURLServer) GetUserURLs(req *pb.GetUserURLsRequest, stream pb.ShortURLService_GetUserURLsServer) error {
	logger.Info(`Used gRPC "GetUserURLs" handler`)

	user := userauth.GetUser(stream.Context())
	if user == nil {
		return status.Error(codes.Unauthenticated, "unauthorized")
	}
	urlsChan, errChan := service.GetAll(stream.Context(), user)

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case err := <-errChan:
			if err != nil {
				var historyNotFoundErr *models.HistoryNotFoundErr
				if errors.As(err, &historyNotFoundErr) {
					return status.Error(codes.NotFound, "no URLs found")
				}
				logger.Error("gRPC GetUserURLs error: ", err)
				return status.Error(codes.Internal, "internal error")
			}
			return nil
		case historyEl, ok := <-urlsChan:
			if !ok {
				return nil
			}
			response := &pb.GetUserURLsResponse{
				OriginalUrl: historyEl.OriginalURL,
				ShortUrl:    historyEl.ShortURL,
			}
			if err := stream.Send(response); err != nil {
				logger.Error("gRPC GetUserURLs stream error: ", err)
				return status.Error(codes.Internal, "stream error")
			}
		}
	}
}

// RemoveURLs удаляет несколько URL пользователя
func (s *ShortURLServer) RemoveURLs(ctx context.Context, req *pb.RemoveURLsRequest) (*emptypb.Empty, error) {
	logger.Info(`Used gRPC "RemoveURLs" handler`)

	if len(req.ShortNames) == 0 {
		return nil, status.Error(codes.InvalidArgument, "short_names required")
	}

	user := userauth.GetUser(ctx)
	if user == nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}
	err := service.RemoveBatch(ctx, user, req.ShortNames)
	if err != nil {
		logger.Error("gRPC RemoveURLs error: ", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &emptypb.Empty{}, nil
}

// Ping проверяет готовность сервиса
func (s *ShortURLServer) Ping(ctx context.Context, req *emptypb.Empty) (*pb.PingResponse, error) {
	logger.Info(`Used gRPC "Ping" handler`)

	ready := service.IsDBReady(ctx)
	return &pb.PingResponse{
		Ready: ready,
	}, nil
}

// GetStats возвращает статистику сервиса
func (s *ShortURLServer) GetStats(ctx context.Context, req *emptypb.Empty) (*pb.GetStatsResponse, error) {
	logger.Info(`Used gRPC "GetStats" handler`)

	stats, err := service.GetStats(ctx)
	if err != nil {
		logger.Error("gRPC GetStats error: ", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.GetStatsResponse{
		Urls:  int32(stats.Urls),
		Users: int32(stats.Users),
	}, nil
}
