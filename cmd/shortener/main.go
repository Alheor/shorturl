package main

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const shortName = `EwHXdJfB`
const addr = `localhost:8080`

var urlMap = make(map[string]string)

func addURL(w http.ResponseWriter, r *http.Request) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `Request body is empty`, http.StatusBadRequest)
		return
	}

	reqURL := strings.TrimSpace(string(reqBody))
	if reqURL == "" {
		http.Error(w, `Request body is empty`, http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(reqURL)
	if err != nil {
		http.Error(w, `Only valid url allowed`, http.StatusBadRequest)
		return
	}

	urlMap[shortName] = reqURL

	w.Header().Add(`Content-Type`, `text/plain; charset=utf-8`)
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(`http://` + addr + `/` + shortName))
	if err != nil {
		panic(err)
	}
}

func getURL(w http.ResponseWriter, r *http.Request) {

	urlID := chi.URLParam(r, "id")
	if urlID == "" {
		http.Error(w, `Identifier is empty`, http.StatusBadRequest)
		return
	}

	location, exists := urlMap[urlID]
	if !exists {
		http.Error(w, `Unknown identifier`, http.StatusBadRequest)
		return
	}

	w.Header().Add(`Location`, location)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	err := http.ListenAndServe(addr, getRouter())
	if err != nil {
		panic(err)
	}
}

func getRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", addURL)
	r.Get("/{id}", getURL)

	return r
}
