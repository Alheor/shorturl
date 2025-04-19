package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Alheor/shorturl/internal/compress"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/httphandler"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/router"
	"github.com/Alheor/shorturl/internal/service"
	"github.com/Alheor/shorturl/internal/urlhasher"
	"github.com/Alheor/shorturl/internal/userauth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type want struct {
	code     int
	response string
	headers  map[string]string
}
type testData struct {
	name        string
	requestBody []byte
	headers     map[string]string
	method      string
	URL         string
	cookie      *http.Cookie
	want        want
}

func TestAddUrl(t *testing.T) {
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
			name:        "generate short url success",
			requestBody: []byte(targetURL + `/test`),
			URL:         `/`,
			method:      http.MethodPost,
			headers:     map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			cookie:      getCookie(),
			want: want{
				code:     http.StatusCreated,
				response: cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`),
				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			},
		},
		{
			name:        "generate short url success with gzip compression",
			requestBody: []byte(targetURL + `/test`),
			URL:         `/`,
			method:      http.MethodPost,
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding:  httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:     httphandler.HeaderContentTypeXGzip,
				httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
			},
			cookie: getCookie(),
			want: want{
				code:     http.StatusConflict,
				response: cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`),
				headers: map[string]string{
					httphandler.HeaderContentType:     httphandler.HeaderContentTypeTextPlain,
					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
				},
			},
		},
		{
			name:        "generate short url with empty body",
			requestBody: []byte(``),
			URL:         `/`,
			method:      http.MethodPost,
			headers:     map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			want: want{
				code:     http.StatusBadRequest,
				response: "URL required\n",
				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			},
		},
	}

	runTests(t, tests)
}

func TestGetUrl(t *testing.T) {
	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = os.Remove(cfg.FileStoragePath)
	require.NoError(t, err)

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	_, err = repository.GetRepository().Add(ctx, user, targetURL+`/test`)
	require.NoError(t, err)

	tests := []testData{
		{
			name:    "get url by short name success",
			URL:     `/` + urlhasher.GetHash(targetURL+`/test`),
			method:  http.MethodGet,
			headers: map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			cookie:  getCookie(),
			want: want{
				code:    http.StatusTemporaryRedirect,
				headers: map[string]string{httphandler.HeaderLocation: targetURL + `/test`},
			},
		},
		{
			name:   "get url by short name success with gzip compression",
			URL:    `/` + urlhasher.GetHash(targetURL+`/test`),
			method: http.MethodGet,
			headers: map[string]string{
				httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain,
			},
			cookie: getCookie(),
			want: want{
				code:    http.StatusTemporaryRedirect,
				headers: map[string]string{httphandler.HeaderLocation: targetURL + `/test`},
			},
		},
		{
			name:    "get url unknown identifier error",
			URL:     `/UnknownIdentifier`,
			method:  http.MethodGet,
			headers: map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			want: want{
				code:     http.StatusBadRequest,
				response: "Unknown identifier\n",
				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			},
		}, {
			name:    "get url empty identifier error",
			URL:     `/`,
			method:  http.MethodGet,
			headers: map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			want: want{
				code:     http.StatusBadRequest,
				response: "Identifier required\n",
				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			},
		},
	}

	runTests(t, tests)
}

func TestGetPing(t *testing.T) {
	cfg := config.Load()

	err := logger.Init(nil)
	require.NoError(t, err)

	httphandler.Init(&cfg)
	service.Init(&cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = repository.Init(ctx, &cfg, nil)
	require.NoError(t, err)

	tests := []testData{
		{
			name:    "get db is ready",
			URL:     `/ping`,
			method:  http.MethodGet,
			headers: map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			cookie:  getCookie(),
			want: want{
				code: http.StatusOK,
			},
		},
	}

	runTests(t, tests)
}

//func TestAddUrlUniqIndexError(t *testing.T) {
//
//	t.Skip(`Run with database only`) // Для ручного запуска с локальной БД
//
//	cfg := config.Load()
//
//	err := logger.Init(nil)
//	require.NoError(t, err)
//
//	httphandler.Init(&cfg)
//	service.Init(&cfg)
//
//	cfg.DatabaseDsn = `user=app password=pass host=localhost port=5432 dbname=app pool_max_conns=10`
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//	defer cancel()
//
//	err = repository.Init(ctx, &cfg, nil)
//	require.NoError(t, err)
//
//	_, err = repository.Connection.Exec(ctx, `TRUNCATE short_url`)
//	require.NoError(t, err)
//
//	_, err = repository.GetRepository().Add(context.Background(), user, targetURL+`/test`)
//	require.NoError(t, err)
//
//	tests := []testData{
//		{
//			name:        "generate short url success",
//			requestBody: []byte(targetURL + `/test`),
//			URL:         `/`,
//			method:      http.MethodPost,
//			headers:     map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			cookie:      getCookie(),
//			want: want{
//				code:     http.StatusConflict,
//				response: cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`),
//				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			},
//		},
//	}
//
//	runTests(t, tests)
//}

func runTests(t *testing.T, tests []testData) {

	ts := httptest.NewServer(router.GetRoutes())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var err error

			if test.headers[httphandler.HeaderContentEncoding] == httphandler.HeaderContentEncodingGzip {
				test.requestBody, err = compress.Compress(test.requestBody)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(test.method, ts.URL+test.URL, bytes.NewReader(test.requestBody))
			require.NoError(t, err)

			if test.cookie != nil {
				req.AddCookie(test.cookie)
			}

			for hName, hVal := range test.headers {
				req.Header.Set(hName, hVal)
			}

			client := ts.Client()
			transport := &http.Transport{Proxy: http.ProxyFromEnvironment}
			transport.DisableCompression = true
			client.Transport = transport

			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.want.code, resp.StatusCode)

			defer resp.Body.Close()
			resBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if test.headers[httphandler.HeaderContentEncoding] == httphandler.HeaderContentEncodingGzip {
				test.requestBody, err = compress.GzipDecompress(resBody)
				require.NoError(t, err)

				assert.Equal(t, test.want.response, string(test.requestBody))
			} else {
				assert.Equal(t, test.want.response, string(resBody))
			}

			for hName, hValue := range test.want.headers {
				assert.Equal(t, hValue, resp.Header.Get(hName))
			}
		})
	}
}

func getCookie() *http.Cookie {
	cookiesValue := string(userauth.GetSignature(user.ID)) + user.ID

	return &http.Cookie{
		Name:  models.CookiesName,
		Value: base64.StdEncoding.EncodeToString([]byte(cookiesValue)),
	}
}
