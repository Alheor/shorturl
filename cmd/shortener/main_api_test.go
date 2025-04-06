package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Alheor/shorturl/internal/compress"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/httphandler"
	"github.com/Alheor/shorturl/internal/logger"
	"github.com/Alheor/shorturl/internal/models"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/urlhasher"
	"github.com/Alheor/shorturl/internal/urlhasher/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const targetURL = `https://practicum.yandex.ru`

func TestApiAddUrl(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = repository.Init(ctx, nil)
	require.NoError(t, err)

	mockRepo := new(mocks.MockShortName)
	mockRepo.On("Generate").Return(`mockStr`)
	urlhasher.Init(mockRepo)

	tests := []testData{
		{
			name:        `API generate short url success`,
			requestBody: []byte(`{"url":"` + targetURL + `/test"}`),
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten`,
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"` + config.GetOptions().BaseHost + `/` + urlhasher.GetShortNameGenerator().Generate() + `"}`,
				headers: map[string]string{
					httphandler.HeaderContentType:     httphandler.HeaderContentTypeJSON,
					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
				},
			},
		},
		{
			name:        `API generate short url success with application/x-gzip header`,
			requestBody: []byte(`{"url":"` + targetURL + `/test"}`),
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentTypeXGzip,
			},
			method: http.MethodPost,
			URL:    `/api/shorten`,
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"` + config.GetOptions().BaseHost + `/` + urlhasher.GetShortNameGenerator().Generate() + `"}`,
				headers: map[string]string{
					httphandler.HeaderContentType:     httphandler.HeaderContentTypeJSON,
					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.URL, bytes.NewReader(test.requestBody))

			var err error
			if test.headers[httphandler.HeaderContentType] == httphandler.HeaderContentTypeXGzip {
				test.requestBody, err = compress.Compress(test.requestBody)

				require.NoError(t, err)
			}

			for hName, hVal := range test.headers {
				req.Header.Set(hName, hVal)
			}

			resp := httptest.NewRecorder()
			httphandler.AddShorten(resp, req)

			res := resp.Result()

			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.headers[httphandler.HeaderContentType], res.Header.Get(httphandler.HeaderContentType))
		})
	}
}

func TestApiAddBatchUrlsSuccess(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = repository.Init(ctx, nil)
	require.NoError(t, err)

	mockRepo := new(mocks.MockShortName)
	mockRepo.On("Generate").Return(`mockStr`)
	urlhasher.Init(mockRepo)

	tests := []testData{
		{
			name:        `API add batch urls success`,
			requestBody: []byte(`[{"correlation_id": "id1","original_url": "` + targetURL + `/test1"},{"correlation_id": "id2","original_url": "` + targetURL + `/test2"}]`),
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code: http.StatusCreated,
				headers: map[string]string{
					httphandler.HeaderContentType:     httphandler.HeaderContentTypeJSON,
					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.URL, bytes.NewReader(test.requestBody))

			var err error
			if test.headers[httphandler.HeaderContentType] == httphandler.HeaderContentTypeXGzip {
				test.requestBody, err = compress.Compress(test.requestBody)

				require.NoError(t, err)
			}

			for hName, hVal := range test.headers {
				req.Header.Set(hName, hVal)
			}

			resp := httptest.NewRecorder()
			httphandler.AddShortenBatch(resp, req)

			res := resp.Result()

			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)

			var response []models.APIBatchResponseEl
			err = json.Unmarshal(resBody, &response)
			require.NoError(t, err)

			assert.Len(t, response, 2)

			assert.Equal(t, `mockStr`, response[0].ShortURL)
			assert.Equal(t, `mockStr`, response[1].ShortURL)

			assert.True(t, response[0].CorrelationID == `id1` || response[0].CorrelationID == `id2`)
			assert.True(t, response[1].CorrelationID == `id1` || response[1].CorrelationID == `id2`)

			assert.Equal(t, test.want.headers[httphandler.HeaderContentType], res.Header.Get(httphandler.HeaderContentType))
		})
	}
}

func TestApiAddBatchUrlsError(t *testing.T) {
	err := logger.Init(nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = repository.Init(ctx, nil)
	require.NoError(t, err)

	mockRepo := new(mocks.MockShortName)
	mockRepo.On("Generate").Return(`mockStr`)
	urlhasher.Init(mockRepo)

	tests := []testData{
		{
			name:        `API add batch urls fail empty body`,
			requestBody: []byte(``),
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					httphandler.HeaderContentType:     httphandler.HeaderContentTypeJSON,
					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
				},
			},
		}, {
			name:        `API add batch urls fail empty array`,
			requestBody: []byte(`[]`),
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					httphandler.HeaderContentType:     httphandler.HeaderContentTypeJSON,
					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
				},
			},
		}, {
			name:        `API add batch urls fail empty object`,
			requestBody: []byte(`[{}]`),
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					httphandler.HeaderContentType:     httphandler.HeaderContentTypeJSON,
					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
				},
			},
		}, {
			name:        `API add batch urls fail invalid url`,
			requestBody: []byte(`[{"correlation_id": "id1","original_url": "invalid url"}]`),
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					httphandler.HeaderContentType:     httphandler.HeaderContentTypeJSON,
					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
				},
			},
		}, {
			name:        `API add batch urls fail invalid object`,
			requestBody: []byte(`[{"correlation_id": "id1"}]`),
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					httphandler.HeaderContentType:     httphandler.HeaderContentTypeJSON,
					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
				},
			},
		}, {
			name:        `API add batch urls fail invalid object`,
			requestBody: []byte(`[{"original_url": "` + targetURL + `/test1"}]`),
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten/batch`,
			want: want{
				code: http.StatusBadRequest,
				headers: map[string]string{
					httphandler.HeaderContentType:     httphandler.HeaderContentTypeJSON,
					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.URL, bytes.NewReader(test.requestBody))

			var err error
			if test.headers[httphandler.HeaderContentType] == httphandler.HeaderContentTypeXGzip {
				test.requestBody, err = compress.Compress(test.requestBody)

				require.NoError(t, err)
			}

			for hName, hVal := range test.headers {
				req.Header.Set(hName, hVal)
			}

			resp := httptest.NewRecorder()
			httphandler.AddShortenBatch(resp, req)

			res := resp.Result()

			defer res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.headers[httphandler.HeaderContentType], res.Header.Get(httphandler.HeaderContentType))
		})
	}
}
