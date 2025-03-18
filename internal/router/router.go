package router

import (
	"net/http"

	"github.com/Alheor/shorturl/internal/handler"
	"github.com/Alheor/shorturl/internal/logger"

	"github.com/go-chi/chi/v5"
)

type HTTPMiddleware func(f http.HandlerFunc) http.HandlerFunc

// GetRoutes Загрузка маршрутизации
func GetRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get(`/*`, middlewareConveyor(handler.GetURL, logger.LoggingHTTPHandler))
	r.Post(`/`, middlewareConveyor(handler.AddURL, logger.LoggingHTTPHandler))
	r.Post(`/api/shorten`, middlewareConveyor(handler.AddShorten, logger.LoggingHTTPHandler))

	return r
}

func middlewareConveyor(h http.HandlerFunc, middlewares ...HTTPMiddleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}
