// Short url service
package main

import (
	"encoding/json"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/log"
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
	//ErrorEmptyRequestBody error message
	ErrorEmptyRequestBody = `Request body is empty`

	//ErrorInvalidURL error message
	ErrorInvalidURL = `Only valid url allowed`

	//ErrorOnlyJsonDataAllowed error message
	ErrorOnlyJsonDataAllowed = `Only json data allowed`

	//ErrorEmptyUrl error message
	ErrorEmptyUrl = `Url is empty`

	//HeaderContentTypeName header "Content-Type" name
	HeaderContentTypeName = `Content-Type`

	//HeaderContentTypeTextPlainValue header "Content-Type" text/plain
	HeaderContentTypeTextPlainValue = `text/plain; charset=utf-8`

	//HeaderContentTypeJsonValue header "Content-Type" application/json
	HeaderContentTypeJsonValue = `application/json`

	//HeaderLocation header "Location" name
	HeaderLocation = `Location`
)

var (
	randomShortName     = randomname.Init()
	shortNameRepository = repository.Init()
	logger              = log.Init(config.Options.LogLevel)
)

type ApiResponse struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type ApiRequest struct {
	Url string `json:"url"`
}

func addURL(w http.ResponseWriter, r *http.Request) {

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqURL := strings.TrimSpace(string(reqBody))
	if reqURL == "" {
		http.Error(w, ErrorEmptyRequestBody, http.StatusBadRequest)
		return
	}

	_, err = url.ParseRequestURI(reqURL)
	if err != nil {
		http.Error(w, ErrorInvalidURL, http.StatusBadRequest)
		return
	}

	shortName, err := appendUrl(reqURL)
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
		http.Error(w, repository.ErrorIDNotFound, http.StatusBadRequest)
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

	var response ApiResponse

	contentType := r.Header.Get(HeaderContentTypeName)
	if contentType != HeaderContentTypeJsonValue {
		response = ApiResponse{Error: ErrorOnlyJsonDataAllowed}
		sendApiResponse(w, &response, http.StatusBadRequest)

		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		response = ApiResponse{Error: err.Error()}
		sendApiResponse(w, &response, http.StatusInternalServerError)

		return
	}

	var request ApiRequest

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		response = ApiResponse{Error: ErrorOnlyJsonDataAllowed}
		sendApiResponse(w, &response, http.StatusBadRequest)

		return
	}

	reqURL := strings.TrimSpace(request.Url)
	if reqURL == "" {
		response = ApiResponse{Error: ErrorEmptyUrl}
		sendApiResponse(w, &response, http.StatusBadRequest)

		return
	}

	_, err = url.ParseRequestURI(reqURL)
	if err != nil {
		response = ApiResponse{Error: ErrorInvalidURL}
		sendApiResponse(w, &response, http.StatusBadRequest)

		return
	}

	shortName, err := appendUrl(reqURL)
	if err != nil {
		response = ApiResponse{Error: err.Error()}
		sendApiResponse(w, &response, http.StatusBadRequest)

		return
	}

	response = ApiResponse{Result: strings.TrimRight(config.Options.BaseHost, `/`) + `/` + shortName}
	sendApiResponse(w, &response, http.StatusCreated)
}

func sendApiResponse(w http.ResponseWriter, apiResponse *ApiResponse, statusCode int) {

	if statusCode == http.StatusInternalServerError {
		logger.Log.Panic(apiResponse.Error)
	}

	w.Header().Set(HeaderContentTypeName, HeaderContentTypeJsonValue)

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

func appendUrl(reqURL string) (string, error) {
	shortName := randomShortName.Generate()

	//try 1
	err := shortNameRepository.Add(shortName, reqURL)
	if err != nil {
		shortName = randomShortName.Generate()

		//try 2
		err = shortNameRepository.Add(shortName, reqURL)
		if err != nil {
			return ``, err
		}
	}

	return shortName, nil
}

func main() {
	config.Load()

	logger.Log.Info("Starting server", zap.String("addr", config.Options.Addr))

	err := http.ListenAndServe(config.Options.Addr, getRouter())
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func getRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", logger.WithLogging(addURL))
	r.Post("/api/shorten", logger.WithLogging(apiShorten))

	r.Get("/{id}", logger.WithLogging(getURL))

	return r
}
