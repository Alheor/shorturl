package main

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/httphandler"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/service"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/stretchr/testify/require"
)

const targetURL = `https://practicum.yandex.ru`

var user = &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}

func TestApiAddUrl(t *testing.T) {
	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = os.Remove(cfg.FileStoragePath)
	require.NoError(t, err)

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	tests := []testData{
		{
			name:        `API generate short url success`,
			requestBody: []byte(`{"url":"` + targetURL + `/test"}`),
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`) + `"}`,
				headers: map[string]string{
					httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
				},
			},
		},
		{
			name:        `API generate short url success with application/x-gzip header`,
			requestBody: []byte(`{"url":"` + targetURL + `/test"}`),
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusConflict,
				response: `{"result":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`) + `"}`,
				headers: map[string]string{
					httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
				},
			},
		},
		{
			name:        `API generate short with empty body error`,
			requestBody: []byte(``),
			headers:     map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON},
			},
		},
		{
			name:        `API generate short with empty url error`,
			requestBody: []byte(`{"url":""}`),
			headers:     map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON},
			},
		},
		{
			name:        `API generate short without url field error`,
			requestBody: []byte(`{"url_test":""}`),
			headers:     map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON},
			},
		}, {
			name:        `API generate short with empty json doc error`,
			requestBody: []byte(`{}`),
			headers:     map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON},
			},
		},
	}

	runTests(t, tests)
}

func TestApiAddBatchUrlsSuccess(t *testing.T) {
	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	tests := []testData{
		{
			name:        `API add batch urls success`,
			requestBody: []byte(`[{"correlation_id":"id1","original_url": "` + targetURL + `/test1"},{"correlation_id":"id2","original_url":"` + targetURL + `/test2"}]`),
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusCreated,
				response: `[{"correlation_id":"id1","short_url":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test1`) + `"},{"correlation_id":"id2","short_url":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test2`) + `"}]`,
				headers: map[string]string{
					httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
				},
			},
		},
	}

	runTests(t, tests)
}

func TestApiAddAndGetBatchUrlsSuccess(t *testing.T) {
	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	tests := []testData{
		{
			name:        `API add batch urls success`,
			requestBody: []byte(`[{"correlation_id":"id1","original_url":"` + targetURL + `/test1"}]`),
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusCreated,
				response: `[{"correlation_id":"id1","short_url":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test1`) + `"}]`,
				headers: map[string]string{
					httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:    `API get url success`,
			headers: map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			method:  http.MethodGet,
			URL:     `/` + urlhasher.GetHash(targetURL+`/test1`),
			cookie:  getCookie(),
			want: want{
				code:    http.StatusTemporaryRedirect,
				headers: map[string]string{httphandler.HeaderLocation: targetURL + `/test1`},
			},
		},
	}

	runTests(t, tests)
}

func TestApiAddBatchUrlsError(t *testing.T) {
	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	tests := []testData{
		{
			name:        `API add batch urls fail empty body`,
			requestBody: []byte(``),
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"invalid body"}`,
				headers: map[string]string{
					httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:        `API add batch urls fail empty array`,
			requestBody: []byte(`[]`),
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"empty url list"}`,
				headers: map[string]string{
					httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:        `API add batch urls fail empty object`,
			requestBody: []byte(`[{}]`),
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"Url '' invalid"}`,
				headers: map[string]string{
					httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:        `API add batch urls fail invalid url`,
			requestBody: []byte(`[{"correlation_id": "id1","original_url": "invalid_url"}]`),
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"Url 'invalid_url' invalid"}`,
				headers: map[string]string{
					httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:        `API add batch urls fail invalid object`,
			requestBody: []byte(`[{"correlation_id": "id1"}]`),
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"Url '' invalid"}`,
				headers: map[string]string{
					httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:        `API add batch urls fail invalid object`,
			requestBody: []byte(`[{"original_url": "` + targetURL + `/test1"}]`),
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"empty correlation_id"}`,
				headers: map[string]string{
					httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
				},
			},
		},
	}

	runTests(t, tests)
}

func TestApiAddUrlUniqIndexError(t *testing.T) {

	t.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	httphandler.Init(&cfg)
	service.Init(&cfg)

	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	err = repository.GetRepository().RemoveByOriginalURL(context.Background(), user, targetURL+`/test`)
	require.NoError(t, err)

	_, err = repository.GetRepository().Add(context.Background(), user, targetURL+`/test`)
	require.NoError(t, err)

	tests := []testData{
		{
			name:        `API generate short url success`,
			requestBody: []byte(`{"url":"` + targetURL + `/test"}`),
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusConflict,
				response: `{"result":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`) + `"}`,
				headers: map[string]string{
					httphandler.HeaderContentType: httphandler.HeaderContentTypeJSON,
				},
			},
		},
	}

	runTests(t, tests)
}
