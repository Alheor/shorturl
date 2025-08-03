package main

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/http/handler"
	"github.com/Alheor/shorturl/internal/ip"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/service"
	"github.com/Alheor/shorturl/internal/shutdown"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/stretchr/testify/require"
)

const targetURL = `https://practicum.yandex.ru`

var user = &models.User{ID: `6a30af51-b6ac-63ba-9e1c-5da06e1b610e`}

func TestApiAddUrl(t *testing.T) {
	shutdown.Init()
	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	handler.Init(&cfg)
	service.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(cfg.FileStoragePath)

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	tests := []testData{
		{
			name:        `API generate short url success`,
			requestBody: []byte(`{"url":"` + targetURL + `/test"}`),
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`) + `"}`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		},
		{
			name:        `API generate short url success with application/x-gzip header`,
			requestBody: []byte(`{"url":"` + targetURL + `/test"}`),
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusConflict,
				response: `{"result":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`) + `"}`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		},
		{
			name:        `API generate short with empty body error`,
			requestBody: []byte(``),
			headers:     map[string]string{handler.HeaderContentType: handler.HeaderContentTypeJSON},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{handler.HeaderContentType: handler.HeaderContentTypeJSON},
			},
		},
		{
			name:        `API generate short with empty url error`,
			requestBody: []byte(`{"url":""}`),
			headers:     map[string]string{handler.HeaderContentType: handler.HeaderContentTypeJSON},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{handler.HeaderContentType: handler.HeaderContentTypeJSON},
			},
		},
		{
			name:        `API generate short without url field error`,
			requestBody: []byte(`{"url_test":""}`),
			headers:     map[string]string{handler.HeaderContentType: handler.HeaderContentTypeJSON},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{handler.HeaderContentType: handler.HeaderContentTypeJSON},
			},
		}, {
			name:        `API generate short with empty json doc error`,
			requestBody: []byte(`{}`),
			headers:     map[string]string{handler.HeaderContentType: handler.HeaderContentTypeJSON},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{handler.HeaderContentType: handler.HeaderContentTypeJSON},
			},
		},
	}

	runTests(t, tests)
}

func TestApiAddBatchUrlsSuccess(t *testing.T) {
	shutdown.Init()
	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	handler.Init(&cfg)
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
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusCreated,
				response: `[{"correlation_id":"id1","short_url":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test1`) + `"},{"correlation_id":"id2","short_url":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test2`) + `"}]`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		},
	}

	runTests(t, tests)
}

func TestApiAddAndGetBatchUrlsSuccess(t *testing.T) {
	shutdown.Init()
	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	handler.Init(&cfg)
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
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusCreated,
				response: `[{"correlation_id":"id1","short_url":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test1`) + `"}]`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:    `API get url success`,
			headers: map[string]string{handler.HeaderContentType: handler.HeaderContentTypeTextPlain},
			method:  http.MethodGet,
			URL:     `/` + urlhasher.GetHash(targetURL+`/test1`),
			cookie:  getCookie(),
			want: want{
				code:    http.StatusTemporaryRedirect,
				headers: map[string]string{handler.HeaderLocation: targetURL + `/test1`},
			},
		},
	}

	runTests(t, tests)
}

func TestApiAddBatchUrlsError(t *testing.T) {
	shutdown.Init()
	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	handler.Init(&cfg)
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
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"invalid body"}`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:        `API add batch urls fail empty array`,
			requestBody: []byte(`[]`),
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"empty url list"}`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:        `API add batch urls fail empty object`,
			requestBody: []byte(`[{}]`),
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"Url '' invalid"}`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:        `API add batch urls fail invalid url`,
			requestBody: []byte(`[{"correlation_id": "id1","original_url": "invalid_url"}]`),
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"Url 'invalid_url' invalid"}`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:        `API add batch urls fail invalid object`,
			requestBody: []byte(`[{"correlation_id": "id1"}]`),
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"Url '' invalid"}`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		}, {
			name:        `API add batch urls fail invalid object`,
			requestBody: []byte(`[{"original_url": "` + targetURL + `/test1"}]`),
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"empty correlation_id"}`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		},
	}

	runTests(t, tests)
}

func TestApiAddUrlUniqIndexError(t *testing.T) {

	t.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	shutdown.Init()
	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	handler.Init(&cfg)
	service.Init(&cfg)

	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	_, err = repository.Connection.Exec(ctx, `TRUNCATE short_url`)
	require.NoError(t, err)

	_, err = repository.GetRepository().Add(context.Background(), user, targetURL+`/test`)
	require.NoError(t, err)

	tests := []testData{
		{
			name:        `API generate short url success`,
			requestBody: []byte(`{"url":"` + targetURL + `/test"}`),
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusConflict,
				response: `{"result":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`) + `"}`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		},
	}

	runTests(t, tests)
}

func TestApiGetAllUrlsFromDBSuccess(t *testing.T) {

	t.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	shutdown.Init()
	cfg := config.Load()

	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	err := logger.Init(nil)
	require.NoError(t, err)

	handler.Init(&cfg)
	service.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	_, err = repository.Connection.Exec(ctx, `TRUNCATE short_url`)
	require.NoError(t, err)

	_, err = repository.GetRepository().Add(context.Background(), user, targetURL+`/test1`)
	require.NoError(t, err)

	_, err = repository.GetRepository().Add(context.Background(), user, targetURL+`/test2`)
	require.NoError(t, err)

	tests := []testData{
		{
			name: `API get all urls success`,
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodGet,
			URL:    `/api/user/urls`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusOK,
				response: `[{"original_url":"` + targetURL + `/test1","short_url":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test1`) + `"},{"original_url":"` + targetURL + `/test2","short_url":"` + cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test2`) + `"}]`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		},
	}

	runTests(t, tests)
}

func TestApiGetAllUrlsError(t *testing.T) {

	t.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	shutdown.Init()
	cfg := config.Load()

	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	err := logger.Init(nil)
	require.NoError(t, err)

	handler.Init(&cfg)
	service.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	_, err = repository.Connection.Exec(ctx, `TRUNCATE short_url`)
	require.NoError(t, err)

	tests := []testData{
		{
			name: `API get all urls without user`,
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodGet,
			URL:    `/api/user/urls`,
			cookie: &http.Cookie{
				Name:  models.CookiesName,
				Value: `aW52YWxpZF92YWx1ZV9pbnZhbGlkX3ZhbHVlX2ludmFsaWRfdmFsdWUK`,
			},
			want: want{
				code:     http.StatusUnauthorized,
				response: `{"error":"Unauthorized"}`,
				headers: map[string]string{
					handler.HeaderContentType: handler.HeaderContentTypeJSON,
				},
			},
		},
		{
			name: `API get all urls empty list`,
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodGet,
			URL:    `/api/user/urls`,
			cookie: getCookie(),
			want: want{
				code: http.StatusNoContent,
			},
		},
	}

	runTests(t, tests)
}

func TestApiRemoveBatch(t *testing.T) {

	t.Skip(`Run with database only`) // Для ручного запуска с локальной БД

	shutdown.Init()
	cfg := config.Load()

	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`

	err := logger.Init(nil)
	require.NoError(t, err)

	handler.Init(&cfg)
	service.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	_, err = repository.Connection.Exec(ctx, `TRUNCATE short_url`)
	require.NoError(t, err)

	var user1 = &models.User{ID: `0b32aa55-b2af-63ba-9e1c-5da06e1b610e`}

	hash1, err := repository.GetRepository().Add(context.Background(), user, targetURL+`/test1`)
	require.NoError(t, err)

	hash2, err := repository.GetRepository().Add(context.Background(), user, targetURL+`/test2`)
	require.NoError(t, err)

	hash3, err := repository.GetRepository().Add(context.Background(), user1, targetURL+`/test3`)
	require.NoError(t, err)

	tests := []testData{
		{
			name:        `API remove batch urls success`,
			requestBody: []byte(`["` + hash1 + `", "` + hash2 + `", "` + hash3 + `"]`),
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodDelete,
			URL:    `/api/user/urls`,
			cookie: getCookie(),
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name: `get url success`,
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodGet,
			URL:    `/` + hash1,
			cookie: getCookie(),
			want: want{
				code: http.StatusGone,
			},
		},
		{
			name: `get url success`,
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodGet,
			URL:    `/` + hash2,
			cookie: getCookie(),
			want: want{
				code: http.StatusGone,
			},
		},
		{
			name: `get url success`,
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodGet,
			URL:    `/` + hash3,
			cookie: getCookie(),
			want: want{
				code: http.StatusTemporaryRedirect,
			},
		},
	}

	runTests(t, tests)
}

func TestStatsEmptySuccess(t *testing.T) {

	shutdown.Init()
	cfg := config.Load()
	cfg.TrustedSubnet = `192.168.0.0/24`

	err := logger.Init(nil)
	require.NoError(t, err)

	handler.Init(&cfg)
	service.Init(&cfg)
	ip.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(cfg.FileStoragePath)

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	tests := []testData{
		{
			name: `API stats empty success`,
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
				ip.HeaderXRealIP:          `192.168.0.1`,
			},
			method: http.MethodGet,
			URL:    `/api/internal/stats`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusOK,
				response: `{"urls":0,"users":0}`,
			},
		},
	}

	runTests(t, tests)
}

func TestStatsNotEmptySuccess(t *testing.T) {

	shutdown.Init()
	cfg := config.Load()
	cfg.TrustedSubnet = `192.168.0.0/24`

	err := logger.Init(nil)
	require.NoError(t, err)

	handler.Init(&cfg)
	service.Init(&cfg)
	ip.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(cfg.FileStoragePath)

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	var user1 = &models.User{ID: `0b32aa55-b2af-63ba-9e1c-5da06e1b610e`}

	_, err = repository.GetRepository().Add(context.Background(), user, targetURL+`/test1`)
	require.NoError(t, err)

	_, err = repository.GetRepository().Add(context.Background(), user, targetURL+`/test2`)
	require.NoError(t, err)

	_, err = repository.GetRepository().Add(context.Background(), user1, targetURL+`/test3`)
	require.NoError(t, err)

	tests := []testData{
		{
			name: `API stats not empty success`,
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
				ip.HeaderXRealIP:          `192.168.0.1`,
			},
			method: http.MethodGet,
			URL:    `/api/internal/stats`,
			cookie: getCookie(),
			want: want{
				code:     http.StatusOK,
				response: `{"urls":3,"users":2}`,
			},
		},
	}

	runTests(t, tests)
}

func TestStatsErrors(t *testing.T) {

	shutdown.Init()
	cfg := config.Load()
	cfg.TrustedSubnet = `192.168.0.0/24`

	err := logger.Init(nil)
	require.NoError(t, err)

	handler.Init(&cfg)
	service.Init(&cfg)
	ip.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_ = os.Remove(cfg.FileStoragePath)

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	tests := []testData{
		{
			name: `API stats without header error`,
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
			},
			method: http.MethodGet,
			URL:    `/api/internal/stats`,
			cookie: getCookie(),
			want: want{
				code: http.StatusForbidden,
			},
		},
		{
			name: `API stats invalid subnet error`,
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
				ip.HeaderXRealIP:          `192.168.1.1`,
			},
			method: http.MethodGet,
			URL:    `/api/internal/stats`,
			cookie: getCookie(),
			want: want{
				code: http.StatusForbidden,
			},
		},
		{
			name: `API stats invalid ip error`,
			headers: map[string]string{
				handler.HeaderContentType: handler.HeaderContentTypeJSON,
				ip.HeaderXRealIP:          `invalid ip`,
			},
			method: http.MethodGet,
			URL:    `/api/internal/stats`,
			cookie: getCookie(),
			want: want{
				code: http.StatusForbidden,
			},
		},
	}

	runTests(t, tests)
}
