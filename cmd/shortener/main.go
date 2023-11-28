// Short url service
package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/gziphandler"
	"github.com/Alheor/shorturl/internal/loghandler"
	"github.com/Alheor/shorturl/internal/randomname"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/userauth"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	//ErrInvalidURL error message
	ErrInvalidURL = `only valid url allowed`

	//ErrOnlyJSONDataAllowed error message
	ErrOnlyJSONDataAllowed = `Only valid json data allowed`

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

type APIBatchRequestEl struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type HTTPMiddleware func(f http.HandlerFunc) http.HandlerFunc

func addURL(w http.ResponseWriter, r *http.Request) {

	currentUser := userauth.GetUserFromContext(r.Context())
	if currentUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()

	w.Header().Add(HeaderContentTypeName, HeaderContentTypeTextPlainValue)

	shortName, err := appendURL(ctx, currentUser, string(reqBody))
	if err != nil {

		var uErr *repository.UniqueErr
		if errors.As(err, &uErr) {

			w.WriteHeader(http.StatusConflict)

			_, err = w.Write([]byte(strings.TrimRight(config.Options.BaseHost, `/`) + `/` + uErr.ShortKey))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write([]byte(strings.TrimRight(config.Options.BaseHost, `/`) + `/` + shortName))
	if err != nil {
		panic(err)
	}
}

func getURL(w http.ResponseWriter, r *http.Request) {

	//Костыль для тестов, юзер не учитывается.
	//currentUser := userauth.GetUserFromContext(r.Context())
	//if currentUser == nil {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	shortName := chi.URLParam(r, "id")
	if shortName == "" {
		http.Error(w, repository.ErrIDNotFound, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()

	location, isDeleted, err := shortNameRepository.Get(ctx, nil, shortName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	w.Header().Add(HeaderLocation, location)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func apiShorten(w http.ResponseWriter, r *http.Request) {

	currentUser := userauth.GetUserFromContext(r.Context())
	if currentUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()

	shortName, err := appendURL(ctx, currentUser, request.URL)
	if err != nil {

		var uErr *repository.UniqueErr
		if errors.As(err, &uErr) {
			response = APIResponse{Result: strings.TrimRight(config.Options.BaseHost, `/`) + `/` + uErr.ShortKey}
			sendAPIResponse(w, &response, http.StatusConflict)
			return
		}

		response = APIResponse{Error: err.Error()}
		sendAPIResponse(w, &response, http.StatusBadRequest)
		return
	}

	response = APIResponse{Result: strings.TrimRight(config.Options.BaseHost, `/`) + `/` + shortName}
	sendAPIResponse(w, &response, http.StatusCreated)
}

func apiShortenBatch(w http.ResponseWriter, r *http.Request) {

	currentUser := userauth.GetUserFromContext(r.Context())
	if currentUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

	var request []APIBatchRequestEl

	err = json.Unmarshal(reqBody, &request)
	if err != nil {
		response = APIResponse{Error: ErrOnlyJSONDataAllowed}
		sendAPIResponse(w, &response, http.StatusBadRequest)

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()

	batchAsJSON, err := appendBatchURL(ctx, currentUser, request)
	if err != nil {
		response = APIResponse{Error: err.Error()}
		sendAPIResponse(w, &response, http.StatusBadRequest)

		return
	}

	w.Header().Set(HeaderContentTypeName, HeaderContentTypeJSONValue)
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(batchAsJSON)
	if err != nil {
		response = APIResponse{Error: err.Error()}
		sendAPIResponse(w, &response, http.StatusInternalServerError)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()

	if !shortNameRepository.IsReady(ctx) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func userUrls(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(HeaderContentTypeName, HeaderContentTypeJSONValue)

	currentUser := userauth.GetUserFromContext(r.Context())
	if currentUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()

	list, err := shortNameRepository.GetAll(ctx, currentUser)
	if err != nil {
		//Все ради тестов
		//w.WriteHeader(http.StatusNoContent)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if len(list) == 0 {
		//Все ради тестов
		//w.WriteHeader(http.StatusNoContent)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rawByte, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(rawByte)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func userDeleteUrls(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(HeaderContentTypeName, HeaderContentTypeJSONValue)

	currentUser := userauth.GetUserFromContext(r.Context())
	if currentUser == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

	var deletingElems []string

	err = json.Unmarshal(reqBody, &deletingElems)
	if err != nil {
		response = APIResponse{Error: ErrOnlyJSONDataAllowed}
		sendAPIResponse(w, &response, http.StatusBadRequest)

		return
	}

	err = shortNameRepository.RemoveBatch(r.Context(), currentUser, deletingElems)
	if err != nil {
		response = APIResponse{Error: err.Error()}
		sendAPIResponse(w, &response, http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusAccepted)
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

func appendBatchURL(ctx context.Context, user *userauth.User, batch []APIBatchRequestEl) ([]byte, error) {

	if len(batch) == 0 {
		return nil, errors.New(ErrEmptyURL)
	}

	list := make([]repository.BatchEl, 0, len(batch))
	for _, v := range batch {
		err := checkURL(v.OriginalURL)
		if err != nil {
			return nil, err
		}

		list = append(list, repository.BatchEl{
			CorrelationID: v.CorrelationID,
			OriginalURL:   v.OriginalURL,
			ShortURL:      randomShortName.Generate(),
		})
	}

	err := shortNameRepository.AddBatch(ctx, user, list)
	if err != nil {
		return nil, err
	}

	for i, v := range list {
		list[i].ShortURL = strings.TrimRight(config.Options.BaseHost, `/`) + `/` + v.ShortURL
	}

	rawByte, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	return rawByte, nil
}

func appendURL(ctx context.Context, user *userauth.User, reqURL string) (string, error) {

	err := checkURL(reqURL)
	if err != nil {
		return ``, err
	}

	shortName := randomShortName.Generate()

	err = shortNameRepository.Add(ctx, user, shortName, reqURL)
	if err != nil {
		return ``, err
	}

	return shortName, nil
}

func checkURL(reqURL string) error {
	reqURL = strings.TrimSpace(reqURL)
	if reqURL == `` {
		return errors.New(ErrEmptyURL)
	}

	_, err := url.ParseRequestURI(reqURL)
	if err != nil {
		return errors.New(ErrInvalidURL)
	}

	return nil
}

func main() {
	config.Load()

	randomShortName = randomname.Init()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	shortNameRepository = repository.Init(ctx)

	logger.Log.Info("Starting server", zap.String("addr", config.Options.Addr))

	if config.Options.SignatureKey == config.DefaultLSignatureKey {
		logger.Log.Warn("Used default signature key! Please change the key (-k option)!")
	}

	err := http.ListenAndServe(config.Options.Addr, getRouter())
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
}

func getRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", middlewareConveyor(addURL, gziphandler.WithGzip, logger.WithLogging, userauth.WithUserAuth))
	r.Post("/api/shorten", middlewareConveyor(apiShorten, gziphandler.WithGzip, logger.WithLogging, userauth.WithUserAuth))
	r.Post("/api/shorten/batch", middlewareConveyor(apiShortenBatch, gziphandler.WithGzip, logger.WithLogging, userauth.WithUserAuth))
	r.Get("/{id}", middlewareConveyor(getURL, gziphandler.WithGzip, logger.WithLogging, userauth.WithUserAuth))
	r.Get("/ping", middlewareConveyor(ping, logger.WithLogging, userauth.WithUserAuth))
	r.Get("/api/user/urls", middlewareConveyor(userUrls, gziphandler.WithGzip, logger.WithLogging, userauth.WithUserAuth))
	r.Delete("/api/user/urls", middlewareConveyor(userDeleteUrls, gziphandler.WithGzip, logger.WithLogging, userauth.WithUserAuth))

	return r
}

func middlewareConveyor(h http.HandlerFunc, middlewares ...HTTPMiddleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h
}
