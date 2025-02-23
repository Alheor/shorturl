package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/controller"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/urlhasher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type want struct {
	code        int
	response    string
	Location    string
	contentType string
}

type mockShortNameGenerator struct{}

func (rg mockShortNameGenerator) Generate() string {
	return `mockStr`
}

func TestAddUrlSuccess(t *testing.T) {

	config.Load()
	repository.Init()
	urlhasher.ShortNameGenerator = new(mockShortNameGenerator)

	tests := []struct {
		name string
		want want
	}{
		{
			name: "generate short url success",
			want: want{
				code:        201,
				response:    config.Options.BaseHost + `/` + urlhasher.ShortNameGenerator.Generate(),
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(`https://practicum.yandex.ru/test`)))
			// создаём новый Recorder
			resp := httptest.NewRecorder()
			controller.AddURL(resp, req)

			res := resp.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestAddUrlWithEmptyBody(t *testing.T) {

	repository.Init()
	urlhasher.ShortNameGenerator = new(mockShortNameGenerator)

	tests := []struct {
		name string
		want want
	}{
		{
			name: "generate short url with empty body",
			want: want{
				code:        400,
				response:    "URL required\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			// создаём новый Recorder
			resp := httptest.NewRecorder()
			controller.AddURL(resp, req)

			res := resp.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestAddUrlWithEmptyUrl(t *testing.T) {

	repository.Init()
	urlhasher.ShortNameGenerator = new(mockShortNameGenerator)

	tests := []struct {
		name string
		want want
	}{
		{
			name: "generate short url with empty url",
			want: want{
				code:        400,
				response:    "URL required\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(``)))
			// создаём новый Recorder
			resp := httptest.NewRecorder()
			controller.AddURL(resp, req)

			res := resp.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestGetUrlSuccess(t *testing.T) {

	repository.Init()
	urlhasher.ShortNameGenerator = new(mockShortNameGenerator)

	repository.Add(`https://practicum.yandex.ru/test`)

	tests := []struct {
		name string
		want want
	}{
		{
			name: "get url by short name success",
			want: want{
				code:        307,
				Location:    `https://practicum.yandex.ru/test`,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+urlhasher.ShortNameGenerator.Generate(), nil)
			// создаём новый Recorder
			resp := httptest.NewRecorder()
			controller.GetURL(resp, request)

			res := resp.Result()

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.Location, res.Header.Get(`Location`))
		})
	}
}

func TestGetUrlUnknownIdentifier(t *testing.T) {

	repository.Init()
	urlhasher.ShortNameGenerator = new(mockShortNameGenerator)

	repository.Add(`https://practicum.yandex.ru/test`)

	tests := []struct {
		name string
		want want
	}{
		{
			name: "generate short url success",
			want: want{
				code:        400,
				response:    "Unknown identifier\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/UnknownIdentifier", nil)
			// создаём новый Recorder
			resp := httptest.NewRecorder()
			controller.GetURL(resp, request)

			res := resp.Result()

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.Location, res.Header.Get(`Location`))
		})
	}
}

func TestGetUrlEmptyIdentifier(t *testing.T) {

	repository.Init()
	urlhasher.ShortNameGenerator = new(mockShortNameGenerator)

	repository.Add(`https://practicum.yandex.ru/test`)

	tests := []struct {
		name string
		want want
	}{
		{
			name: "generate short url success",
			want: want{
				code:        400,
				response:    "Identifier required\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			// создаём новый Recorder
			resp := httptest.NewRecorder()
			controller.GetURL(resp, request)

			res := resp.Result()

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.Location, res.Header.Get(`Location`))
		})
	}
}
