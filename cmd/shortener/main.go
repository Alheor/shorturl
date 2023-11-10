// Short url service
package main

import (
	"encoding/json"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/gziphandler"
	"github.com/Alheor/shorturl/internal/loghandler"
	"github.com/Alheor/shorturl/internal/randomname"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	//ErrEmptyRequestBody error message
	ErrEmptyRequestBody = `Request body is empty`

	//ErrInvalidURL error message
	ErrInvalidURL = `Only valid url allowed`

	//ErrOnlyJSONDataAllowed error message
	ErrOnlyJSONDataAllowed = `Only json data allowed`

	//ErrEmptyURL error message
	ErrEmptyURL = `URL is empty`

	//HeaderContentTypeName header "Content-Type" name
	HeaderContentTypeName = `Content-Type`

	//HeaderAcceptEncodingName  header "Accept-Encoding" name
	HeaderAcceptEncodingName = `Accept-Encoding`

	//HeaderContentEncodingName  header "Content-Encoding" name
	HeaderContentEncodingName = `Content-Encoding`

	//HeaderAcceptEncodingValue header "Accept-Encoding" value
	HeaderAcceptEncodingValue = `gzip`

	//HeaderContentTypeTextPlainValue header "Content-Type" text/plain
	HeaderContentTypeTextPlainValue = `text/plain; charset=utf-8`

	//HeaderContentTypeTextHTMLValue header "Content-Type" text/html
	HeaderContentTypeTextHTMLValue = `text/html; charset=utf-8`

	//HeaderContentTypeJSONValue header "Content-Type" application/json
	HeaderContentTypeJSONValue = `application/json`

	//HeaderContentTypeXgzipValue header "Content-Type" application/x-gzip
	HeaderContentTypeXgzipValue = `application/x-gzip`

	//HeaderLocation header "Location" name
	HeaderLocation = `Location`
)

var (
	randomShortName     randomname.RandomStringGenerator
	shortNameRepository repository.Repository

	logger = loghandler.Init()
)

type APIResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type APIRequest struct {
	URL string `json:"url"`
}

type HTTPMiddleware func(f http.HandlerFunc) http.HandlerFunc

func addURL(w http.ResponseWriter, r *http.Request) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqURL := strings.TrimSpace(string(reqBody))
	if reqURL == "" {
		http.Error(w, ErrEmptyRequestBody, http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(reqURL)
	if err != nil {
		http.Error(w, ErrInvalidURL, http.StatusBadRequest)
		return
	}

	shortName, err := appendURL(reqURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add(HeaderContentTypeName, HeaderContentTypeTextPlainValue)
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(strings.TrimRight(config.Options.BaseHost, `/`) + `/` + shortName))
	if err != nil {
		panic(err)
	}
}

func getURL(w http.ResponseWriter, r *http.Request) {

	shortName := chi.URLParam(r, "id")
	if shortName == "" {
		http.Error(w, repository.ErrIDNotFound, http.StatusBadRequest)
		return
	}

	location, err := shortNameRepository.Get(shortName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add(HeaderLocation, location)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func apiShorten(w http.ResponseWriter, r *http.Request) {

	var response APIResponse

	contentType := r.Header.Get(HeaderContentTypeName)
	if contentType != HeaderContentTypeJSONValue && contentType != HeaderContentTypeXgzipValue {
		response = APIResponse{Error: ErrOnlyJSONDataAllowed}
		sendAPIResponse(w, &response, http.StatusBadRequest)

		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		response = APIResponse{Error: err.Error()}
		sendAPIResponse(w, &response, http.StatusInternalServerError)

		return
	}

	var request APIRequest

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		response = APIResponse{Error: ErrOnlyJSONDataAllowed}
		sendAPIResponse(w, &response, http.StatusBadRequest)

		return
	}

	reqURL := strings.TrimSpace(request.URL)
	if reqURL == "" {
		response = APIResponse{Error: ErrEmptyURL}
		sendAPIResponse(w, &response, http.StatusBadRequest)

		return
	}

	_, err = url.ParseRequestURI(reqURL)
	if err != nil {
		response = APIResponse{Error: ErrInvalidURL}
		sendAPIResponse(w, &response, http.StatusBadRequest)

		return
	}

	shortName, err := appendURL(reqURL)
	if err != nil {
		response = APIResponse{Error: err.Error()}
		sendAPIResponse(w, &response, http.StatusBadRequest)

		return
	}

	response = APIResponse{Result: strings.TrimRight(config.Options.BaseHost, `/`) + `/` + shortName}
	sendAPIResponse(w, &response, http.StatusCreated)
}

func ping(w http.ResponseWriter, r *http.Request) {

	if !shortNameRepository.StorageIsReady() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func sendAPIResponse(w http.ResponseWriter, apiResponse *APIResponse, statusCode int) {

	if statusCode == http.StatusInternalServerError {
		logger.Log.Panic(apiResponse.Error)
	}

	w.Header().Set(HeaderContentTypeName, HeaderContentTypeJSONValue)

	rawByte, err := json.Marshal(apiResponse)
	if err != nil {
		logger.Log.Panic(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(statusCode)

	_, err = w.Write(rawByte)
	if err != nil {
		logger.Log.Panic(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func appendURL(reqURL string) (string, error) {
	shortName := randomShortName.Generate()

	err := shortNameRepository.Add(shortName, reqURL)
	if err != nil {
		return ``, err
	}

	return shortName, nil
}

func main() {
	config.Load()

	randomShortName = randomname.Init()
	shortNameRepository = repository.Init()

	logger.Log.Info("Starting server", zap.String("addr", config.Options.Addr))

	err := http.ListenAndServe(config.Options.Addr, getRouter())
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func getRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", middlewareConveyor(addURL, gziphandler.WithGzip, logger.WithLogging))
	r.Post("/api/shorten", middlewareConveyor(apiShorten, gziphandler.WithGzip, logger.WithLogging))
	r.Get("/{id}", middlewareConveyor(getURL, gziphandler.WithGzip, logger.WithLogging))
	r.Get("/ping", middlewareConveyor(ping, logger.WithLogging))

	return r
}

func middlewareConveyor(h http.HandlerFunc, middlewares ...HTTPMiddleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}
