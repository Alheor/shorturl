package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

const shortName = `EwHXdJfB`
const addr = `localhost:8080`

var urlMap = map[string]string{shortName: "http://fj6dgd0jd.yandex/hwlxqpmtr"}

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

	w.Header().Add(`Content-Type`, `text/plain; charset=utf-8`)
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(`http://` + addr + `/` + shortName))
	if err != nil {
		panic(err)
	}
}

func getURL(w http.ResponseWriter, r *http.Request) {

	urlID := strings.TrimSpace(r.RequestURI)
	if urlID == "" {
		http.Error(w, `Invalid url`, http.StatusBadRequest)
		return
	}

	urlID = strings.TrimLeft(urlID, `/`)

	location, exists := urlMap[urlID]
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

		if r.Method == http.MethodPost {
			addURL(w, r)
		} else {
			getURL(w, r)
		}
	})

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		panic(err)
	}
}
