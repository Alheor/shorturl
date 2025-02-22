package controller

import (
	"net/http"
)

const Addr = `localhost:8080`
const Schema = `http://`

// GetRouter Загрузка маршрутизации
func GetRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, func(resp http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			GetURL(resp, req)
			return
		}

		if req.Method == http.MethodPost {
			AddURL(resp, req)
			return
		}

		resp.WriteHeader(http.StatusMethodNotAllowed)
	})

	return mux
}
