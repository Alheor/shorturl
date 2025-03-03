package router

import (
	"github.com/Alheor/shorturl/internal/controller"

	"github.com/go-chi/chi/v5"
)

// GetRoutes Загрузка маршрутизации
func GetRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get(`/*`, controller.GetURL)
	r.Post(`/`, controller.AddURL)

	return r
}
