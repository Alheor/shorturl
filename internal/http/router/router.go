// Package router - сервис маршрутизации.
//
// # Описание
//
// Описывает маршрутизацию и позволяет загрузить ее в веб-сервер.
package router

import (
	"net/http"

	"github.com/Alheor/shorturl/internal/compress"
	"github.com/Alheor/shorturl/internal/http/handler"
	"github.com/Alheor/shorturl/internal/ip"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/userauth"

	"github.com/go-chi/chi/v5"
)

// HTTPMiddleware функция-обертка для реализации конвейера
type HTTPMiddleware func(f http.HandlerFunc) http.HandlerFunc

// GetRoutes Загрузка маршрутизации
func GetRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get(`/*`,
		middlewareConveyor(handler.GetURL, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Get(`/ping`,
		middlewareConveyor(handler.Ping, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Get(`/api/user/urls`,
		middlewareConveyor(handler.GetAllShorten, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Get(`/api/internal/stats`,
		middlewareConveyor(handler.Stats, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, ip.SubnetHTTPHandler))

	r.Delete(`/api/user/urls`,
		middlewareConveyor(handler.DeleteShorten, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Post(`/`,
		middlewareConveyor(handler.AddURL, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Post(`/api/shorten`,
		middlewareConveyor(handler.AddShorten, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Post(`/api/shorten/batch`,
		middlewareConveyor(handler.AddShortenBatch, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	return r
}

// Функция - конвейер.
func middlewareConveyor(h http.HandlerFunc, middlewares ...HTTPMiddleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}
