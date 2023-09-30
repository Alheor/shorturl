package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

const shortName = `EwHXdJfB`
const addr = `localhost:8080`

var urlMap = map[string]string{shortName: "https://practicum.yandex.ru/"}

func addUrl(w http.ResponseWriter, r *http.Request) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `Request body is empty`, http.StatusBadRequest)
		return
	}

	reqUrl := strings.TrimSpace(string(reqBody))
	if reqUrl == "" {
		http.Error(w, `Request body is empty`, http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(reqUrl)
	if err != nil {
		http.Error(w, `Only valid url allowed`, http.StatusBadRequest)
		return
	}

	w.Header().Add(`Content-Type`, `text/plain; charset=utf-8`)
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(`http://` + addr + `/` + shortName))
	if err != nil {
		panic(err)
	}
}

func getUrl(w http.ResponseWriter, r *http.Request) {

	urlId := strings.TrimSpace(r.RequestURI)
	if urlId == "" {
		http.Error(w, `Invalid url`, http.StatusBadRequest)
		return
	}

	urlId = strings.TrimLeft(urlId, `/`)

	location, exists := urlMap[urlId]
	if !exists {
		http.Error(w, `Unknown identifier`, http.StatusBadRequest)
		return
	}

	w.Header().Add(`Location`, location)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get(`Content-Type`) != `text/plain` {
			http.Error(w, `Only text/plain are allowed!`, http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {
			addUrl(w, r)
		} else {
			getUrl(w, r)
		}
	})

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		panic(err)
	}
}
