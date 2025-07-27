package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/service"
	"github.com/Alheor/shorturl/internal/userauth"
)

// HTTP заголовки
const (
	// HeaderContentType header "Content-Type" name.
	HeaderContentType = `Content-Type`

	// HeaderContentEncoding header "Content-Encoding" name.
	HeaderContentEncoding = `Content-Encoding`

	// HeaderAcceptEncoding header "Accept-Encoding" name.
	HeaderAcceptEncoding = `Accept-Encoding`

	// HeaderLocation header "Location" name.
	HeaderLocation = `Location`

	// HeaderContentTypeJSON header Content-Type value application/json.
	HeaderContentTypeJSON = `application/json`

	// HeaderContentTypeXGzip header Content-Type value application/x-gzip.
	HeaderContentTypeXGzip = `application/x-gzip`

	// HeaderContentTypeTextPlain header Content-Type value text/plain.
	HeaderContentTypeTextPlain = `text/plain; charset=utf-8`

	// HeaderContentTypeTextHTML header Content-Type value text/html.
	HeaderContentTypeTextHTML = `text/html`

	// HeaderContentEncodingGzip header Content-Encoding value gzip.
	HeaderContentEncodingGzip = `gzip`
)

var baseHost string

// Init Подготовка HTTP обработчиков к работе.
func Init(config *config.Options) {
	baseHost = config.BaseHost
}

// AddURL Обработчик запроса на добавление URL пользователя.
func AddURL(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "AddURL" handler`)

	var body []byte
	var err error

	if body, err = io.ReadAll(req.Body); err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	URL := strings.TrimSpace(string(body))
	if len(URL) == 0 {
		http.Error(resp, `URL required`, http.StatusBadRequest)
		return
	}

	if _, err = url.ParseRequestURI(URL); err != nil {
		http.Error(resp, `Only valid url allowed`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	user := userauth.GetUser(ctx)
	if user == nil {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp.Header().Add(HeaderContentType, HeaderContentTypeTextPlain)

	shortURL, err := service.Add(ctx, user, URL)
	if err != nil {

		var uniqErr *models.UniqueErr
		if errors.As(err, &uniqErr) {

			resp.WriteHeader(http.StatusConflict)

			_, err = resp.Write([]byte(baseHost + `/` + uniqErr.ShortKey))
			if err != nil {
				logger.Error(`error while response write`, err)
				resp.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusCreated)

	_, err = resp.Write([]byte(baseHost + `/` + shortURL))
	if err != nil {
		logger.Error(`error while response write`, err)
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

// GetURL Обработчик запроса на получение одного URL пользователя.
func GetURL(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "GetURL" handler`)

	shortName := strings.TrimLeft(strings.TrimSpace(req.RequestURI), `/`)
	if len(shortName) == 0 {
		http.Error(resp, `Identifier required`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	URL, isRemoved := service.Get(ctx, nil, shortName)
	if len(URL) == 0 {
		http.Error(resp, `Unknown identifier`, http.StatusBadRequest)
		return
	}

	if isRemoved {
		resp.WriteHeader(http.StatusGone)
		return
	}

	resp.Header().Set(HeaderLocation, URL)
	resp.WriteHeader(http.StatusTemporaryRedirect)
}

// Ping Обработчик запроса на проверку работоспособности сервиса.
func Ping(resp http.ResponseWriter, req *http.Request) {

	logger.Info(`Used "Ping" handler`)

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	if service.IsDBReady(ctx) {
		resp.WriteHeader(http.StatusOK)
		return
	}

	resp.WriteHeader(http.StatusInternalServerError)
}
