// Short url service
package main

import (
	"fmt"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/randomname"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	//ErrorEmptyRequestBody error message
	ErrorEmptyRequestBody = `Request body is empty`

	//ErrorInvalidUrl error message
	ErrorInvalidUrl = `Only valid url allowed`

	//HeaderContentTypeName header "Content-Type" name
	HeaderContentTypeName = `Content-Type`

	//HeaderContentTypeValue header "Content-Type" value
	HeaderContentTypeValue = `text/plain; charset=utf-8`

	//HeaderLocation header "Location" name
	HeaderLocation = `Location`
)

var (
	shortNameRepository repository.Repository
	randomShortName     randomname.RandomString
)

func init() {
	randomShortName = new(randomname.ShortName)
	shortNameRepository = new(repository.ShortName).Init()
}

func addURL(w http.ResponseWriter, r *http.Request) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `Read body error:`+err.Error(), http.StatusBadRequest)
		return
	}

	reqURL := strings.TrimSpace(string(reqBody))
	if reqURL == "" {
		http.Error(w, ErrorEmptyRequestBody, http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(reqURL)
	if err != nil {
		http.Error(w, ErrorInvalidUrl, http.StatusBadRequest)
		return
	}

	shortName := randomShortName.Generate()

	//try 1
	err = shortNameRepository.AddUrl(shortName, reqURL)
	if err != nil {
		shortName = randomShortName.Generate()

		//try 2
		err = shortNameRepository.AddUrl(shortName, reqURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.Header().Add(HeaderContentTypeName, HeaderContentTypeValue)
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(strings.TrimRight(config.Options.BaseHost, `/`) + `/` + shortName))
	if err != nil {
		panic(err)
	}
}

func getURL(w http.ResponseWriter, r *http.Request) {

	shortName := chi.URLParam(r, "id")
	if shortName == "" {
		http.Error(w, repository.ErrorUrlNotFound, http.StatusBadRequest)
		return
	}

	location, err := shortNameRepository.GetUrl(shortName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add(HeaderLocation, location)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	config.Load()

	fmt.Println(`Server listen ` + config.Options.Addr)

	err := http.ListenAndServe(config.Options.Addr, getRouter())
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
