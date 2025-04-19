package main

//
//import (
//	"bytes"
//	"context"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"testing"
//	"time"
//
//	"github.com/Alheor/shorturl/internal/compress"
//	"github.com/Alheor/shorturl/internal/config"
//	"github.com/Alheor/shorturl/internal/httphandler"
//	"github.com/Alheor/shorturl/internal/logger"
//	"github.com/Alheor/shorturl/internal/repository"
//	"github.com/Alheor/shorturl/internal/service"
//	"github.com/Alheor/shorturl/internal/urlhasher"
//
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//)
//
//type want struct {
//	code     int
//	response string
//	headers  map[string]string
//}
//type testData struct {
//	name        string
//	requestBody []byte
//	headers     map[string]string
//	method      string
//	URL         string
//	want        want
//}
//
//func TestAddUrl(t *testing.T) {
//	cfg := config.Load()
//
//	err := logger.Init(nil)
//	require.NoError(t, err)
//
//	httphandler.Init(&cfg)
//	service.Init(&cfg)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
//	defer cancel()
//
//	err = os.Remove(cfg.FileStoragePath)
//	require.NoError(t, err)
//
//	err = repository.Init(ctx, &cfg, nil)
//	require.NoError(t, err)
//
//	tests := []testData{
//		{
//			name:        "generate short url success",
//			requestBody: []byte(targetURL + `/test`),
//			URL:         `/`,
//			method:      http.MethodPost,
//			headers:     map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			want: want{
//				code:     http.StatusCreated,
//				response: cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`),
//				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			},
//		},
//		{
//			name:        "generate short url success with gzip compression",
//			requestBody: []byte(targetURL + `/test`),
//			URL:         `/`,
//			method:      http.MethodPost,
//			headers: map[string]string{
//				httphandler.HeaderAcceptEncoding:  httphandler.HeaderContentEncodingGzip,
//				httphandler.HeaderContentType:     httphandler.HeaderContentTypeXGzip,
//				httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
//			},
//			want: want{
//				code:     http.StatusConflict,
//				response: cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`),
//				headers: map[string]string{
//					httphandler.HeaderContentType:     httphandler.HeaderContentTypeTextPlain,
//					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
//				},
//			},
//		},
//		{
//			name:        "generate short url with empty body",
//			requestBody: []byte(``),
//			URL:         `/`,
//			method:      http.MethodPost,
//			headers:     map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			want: want{
//				code:     http.StatusBadRequest,
//				response: "URL required\n",
//				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			req := httptest.NewRequest(test.method, test.URL, bytes.NewReader(test.requestBody))
//
//			var err error
//			if test.headers[httphandler.HeaderContentType] == httphandler.HeaderContentTypeXGzip {
//				test.requestBody, err = compress.Compress(test.requestBody)
//
//				require.NoError(t, err)
//			}
//
//			for hName, hVal := range test.headers {
//				req.Header.Set(hName, hVal)
//			}
//
//			resp := httptest.NewRecorder()
//			httphandler.AddURL(resp, req)
//
//			res := resp.Result()
//
//			assert.Equal(t, test.want.code, res.StatusCode)
//
//			defer res.Body.Close()
//			resBody, err := io.ReadAll(res.Body)
//
//			require.NoError(t, err)
//			assert.Equal(t, test.want.response, string(resBody))
//			assert.Equal(t, test.want.headers[httphandler.HeaderContentType], res.Header.Get(httphandler.HeaderContentType))
//		})
//	}
//}
//
//func TestGetUrl(t *testing.T) {
//	cfg := config.Load()
//
//	err := logger.Init(nil)
//	require.NoError(t, err)
//
//	httphandler.Init(&cfg)
//	service.Init(&cfg)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//	defer cancel()
//
//	err = os.Remove(cfg.FileStoragePath)
//	require.NoError(t, err)
//
//	err = repository.Init(ctx, &cfg, nil)
//	require.NoError(t, err)
//
//	_, err = repository.GetRepository().Add(ctx, targetURL+`/test`)
//	require.NoError(t, err)
//
//	tests := []testData{
//		{
//			name:    "get url by short name success",
//			URL:     `/` + urlhasher.GetHash(targetURL+`/test`),
//			method:  http.MethodGet,
//			headers: map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			want: want{
//				code:    http.StatusTemporaryRedirect,
//				headers: map[string]string{httphandler.HeaderLocation: targetURL + `/test`},
//			},
//		},
//		{
//			name:   "get url by short name success with gzip compression",
//			URL:    `/` + urlhasher.GetHash(targetURL+`/test`),
//			method: http.MethodGet,
//			headers: map[string]string{
//				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
//				httphandler.HeaderContentType:    httphandler.HeaderContentTypeTextPlain,
//			},
//			want: want{
//				code:    http.StatusTemporaryRedirect,
//				headers: map[string]string{httphandler.HeaderLocation: targetURL + `/test`},
//			},
//		},
//		{
//			name:    "get url unknown identifier error",
//			URL:     `/UnknownIdentifier`,
//			method:  http.MethodGet,
//			headers: map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			want: want{
//				code:     http.StatusBadRequest,
//				response: "Unknown identifier\n",
//				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			},
//		}, {
//			name:    "get url empty identifier error",
//			URL:     `/`,
//			method:  http.MethodGet,
//			headers: map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			want: want{
//				code:     http.StatusBadRequest,
//				response: "Identifier required\n",
//				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			req := httptest.NewRequest(test.method, test.URL, nil)
//
//			for hName, hVal := range test.headers {
//				req.Header.Set(hName, hVal)
//			}
//
//			resp := httptest.NewRecorder()
//			httphandler.GetURL(resp, req)
//
//			res := resp.Result()
//
//			// проверяем код ответа
//			assert.Equal(t, test.want.code, res.StatusCode)
//
//			// получаем и проверяем тело запроса
//			defer res.Body.Close()
//			respBody, err := io.ReadAll(res.Body)
//
//			require.NoError(t, err)
//
//			assert.Equal(t, test.want.response, string(respBody))
//			assert.Equal(t, test.want.headers[httphandler.HeaderContentType], res.Header.Get(httphandler.HeaderContentType))
//		})
//	}
//}
//
//func TestGetPing(t *testing.T) {
//	cfg := config.Load()
//
//	err := logger.Init(nil)
//	require.NoError(t, err)
//
//	httphandler.Init(&cfg)
//	service.Init(&cfg)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
//	defer cancel()
//
//	err = repository.Init(ctx, &cfg, nil)
//	require.NoError(t, err)
//
//	tests := []testData{
//		{
//			name:    "get db is ready",
//			URL:     `/ping`,
//			method:  http.MethodGet,
//			headers: map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			want: want{
//				code: http.StatusOK,
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			req := httptest.NewRequest(test.method, test.URL, nil)
//
//			for hName, hVal := range test.headers {
//				req.Header.Set(hName, hVal)
//			}
//
//			resp := httptest.NewRecorder()
//			httphandler.Ping(resp, req)
//
//			res := resp.Result()
//
//			// проверяем код ответа
//			assert.Equal(t, test.want.code, res.StatusCode)
//
//			// получаем и проверяем тело запроса
//			defer res.Body.Close()
//			respBody, err := io.ReadAll(res.Body)
//
//			require.NoError(t, err)
//
//			assert.Equal(t, test.want.response, string(respBody))
//			assert.Equal(t, test.want.headers[httphandler.HeaderContentType], res.Header.Get(httphandler.HeaderContentType))
//		})
//	}
//}
//
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
//	err = repository.GetRepository().RemoveByOriginalURL(context.Background(), targetURL+`/test`)
//	require.NoError(t, err)
//
//	_, err = repository.GetRepository().Add(context.Background(), targetURL+`/test`)
//	require.NoError(t, err)
//
//	tests := []testData{
//		{
//			name:        "generate short url success",
//			requestBody: []byte(targetURL + `/test`),
//			URL:         `/`,
//			method:      http.MethodPost,
//			headers:     map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			want: want{
//				code:     http.StatusConflict,
//				response: cfg.BaseHost + `/` + urlhasher.GetHash(targetURL+`/test`),
//				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			req := httptest.NewRequest(test.method, test.URL, bytes.NewReader(test.requestBody))
//
//			var err error
//			if test.headers[httphandler.HeaderContentType] == httphandler.HeaderContentTypeXGzip {
//				test.requestBody, err = compress.Compress(test.requestBody)
//
//				require.NoError(t, err)
//			}
//
//			for hName, hVal := range test.headers {
//				req.Header.Set(hName, hVal)
//			}
//
//			resp := httptest.NewRecorder()
//			httphandler.AddURL(resp, req)
//
//			res := resp.Result()
//
//			assert.Equal(t, test.want.code, res.StatusCode)
//
//			defer res.Body.Close()
//			resBody, err := io.ReadAll(res.Body)
//
//			require.NoError(t, err)
//			assert.Equal(t, test.want.response, string(resBody))
//			assert.Equal(t, test.want.headers[httphandler.HeaderContentType], res.Header.Get(httphandler.HeaderContentType))
//		})
//	}
//}
