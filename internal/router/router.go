package router

import (
	"net/http"

	"github.com/Alheor/shorturl/internal/compress"
	"github.com/Alheor/shorturl/internal/httphandler"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/userauth"

	"github.com/go-chi/chi/v5"
)

type HTTPMiddleware func(f http.HandlerFunc) http.HandlerFunc

// GetRoutes Загрузка маршрутизации
func GetRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get(`/*`,
		middlewareConveyor(httphandler.GetURL, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Get(`/ping`,
		middlewareConveyor(httphandler.Ping, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Get(`/api/user/urls`,
		middlewareConveyor(httphandler.GetALlShorten, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Post(`/`,
		middlewareConveyor(httphandler.AddURL, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Post(`/api/shorten`,
		middlewareConveyor(httphandler.AddShorten, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	r.Post(`/api/shorten/batch`,
		middlewareConveyor(httphandler.AddShortenBatch, logger.LoggingHTTPHandler, compress.GzipHTTPHandler, userauth.AuthHTTPHandler))

	return r
}

func middlewareConveyor(h http.HandlerFunc, middlewares ...HTTPMiddleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}
