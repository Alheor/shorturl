package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alheor/shorturl/internal/compress"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/httphandler"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiAddUrlSuccess(t *testing.T) {

	urlhasher.ShortNameGenerator = new(mockShortNameGenerator)

	tests := []testData{
		{
			name:        `API generate short url success`,
			requestBody: []byte(`{"url":"https://practicum.yandex.ru/test"}`),
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentTypeJSON,
			},
			method: http.MethodPost,
			URL:    `/api/shorten`,
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"` + config.GetOptions().BaseHost + `/` + urlhasher.ShortNameGenerator.Generate() + `"}`,
				headers: map[string]string{
					httphandler.HeaderContentType:     httphandler.HeaderContentTypeJSON,
					httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
				},
			},
		},
		{
			name:        `API generate short url success with application/x-gzip header`,
			requestBody: []byte(`{"url":"https://practicum.yandex.ru/test"}`),
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentEncodingXGzip,
			},
			method: http.MethodPost,
			URL:    `/api/shorten`,
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"` + config.GetOptions().BaseHost + `/` + urlhasher.ShortNameGenerator.Generate() + `"}`,
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
			if test.headers[httphandler.HeaderContentType] == httphandler.HeaderContentEncodingXGzip {
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
