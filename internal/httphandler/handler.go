package httphandler

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/service"
)

const (
	//HeaderContentType header "Content-Type" name
	HeaderContentType = `Content-Type`

	//HeaderContentEncoding header "Content-Encoding" name
	HeaderContentEncoding = `Content-Encoding`

	//HeaderAcceptEncoding header "Accept-Encoding" name
	HeaderAcceptEncoding = `Accept-Encoding`

	//HeaderLocation header "Location" name
	HeaderLocation = `Location`

	//HeaderContentTypeJSON header Content-Type value application/json
	HeaderContentTypeJSON = `application/json`

	//HeaderContentTypeXGzip header Content-Type value application/x-gzip
	HeaderContentTypeXGzip = `application/x-gzip`

	//HeaderContentTypeTextPlain header Content-Type value text/plain
	HeaderContentTypeTextPlain = `text/plain; charset=utf-8`

	//HeaderContentTypeTextHTML header Content-Type value text/html
	HeaderContentTypeTextHTML = `text/html`

	//HeaderContentEncodingGzip header Content-Encoding value gzip
	HeaderContentEncodingGzip = `gzip`
)

// AddURL контроллер добавления URL
func AddURL(resp http.ResponseWriter, req *http.Request) {

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

	shortURL, err := service.Add(ctx, URL)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.Header().Add(HeaderContentType, HeaderContentTypeTextPlain)
	resp.WriteHeader(http.StatusCreated)

	_, err = resp.Write([]byte(config.GetOptions().BaseHost + `/` + shortURL))
	if err != nil {
		logger.Error(`error while response write`, err)
		resp.WriteHeader(http.StatusInternalServerError)
	}
}

// GetURL контроллер получения URL по короткому имени
func GetURL(resp http.ResponseWriter, req *http.Request) {

	shortName := strings.TrimLeft(strings.TrimSpace(req.RequestURI), `/`)
	if len(shortName) == 0 {
		http.Error(resp, `Identifier required`, http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	URL := service.Get(ctx, shortName)
	if len(URL) == 0 {
		http.Error(resp, `Unknown identifier`, http.StatusBadRequest)
		return
	}

	resp.Header().Set(HeaderLocation, URL)
	resp.WriteHeader(http.StatusTemporaryRedirect)
}

func Ping(resp http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(req.Context(), 1*time.Second)
	defer cancel()

	if service.IsDBReady(ctx) {
		resp.WriteHeader(http.StatusOK)
		return
	}

	resp.WriteHeader(http.StatusInternalServerError)
}
