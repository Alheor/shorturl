package handler

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/repository"
)

const (
	//HeaderContentTypeName header "Content-Type" name
	HeaderContentTypeName = `Content-Type`

	//HeaderLocation header "Location" name
	HeaderLocation = `Location`

	//HeaderContentTypeJSONValue header "Content-Type" application/json
	HeaderContentTypeJSONValue = `application/json; charset=utf-8`

	//HeaderContentTypeTextPlainValue header "Content-Type" text/plain
	HeaderContentTypeTextPlainValue = `text/plain; charset=utf-8`
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

	shortName := repository.GetRepository().Add(URL)

	resp.Header().Add(HeaderContentTypeName, HeaderContentTypeTextPlainValue)
	resp.WriteHeader(http.StatusCreated)

	_, err = resp.Write([]byte(config.GetOptions().BaseHost + `/` + shortName))
	if err != nil {
		panic(err)
	}
}

// GetURL контроллер получения URL по короткому имени
func GetURL(resp http.ResponseWriter, req *http.Request) {

	shortName := strings.TrimLeft(strings.TrimSpace(req.RequestURI), `/`)
	if len(shortName) == 0 {
		http.Error(resp, `Identifier required`, http.StatusBadRequest)
		return
	}

	URL := repository.GetRepository().GetByShortName(shortName)
	if URL == nil {
		http.Error(resp, `Unknown identifier`, http.StatusBadRequest)
		return
	}

	resp.Header().Set(HeaderLocation, *URL)
	resp.WriteHeader(http.StatusTemporaryRedirect)
}
