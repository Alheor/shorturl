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
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/urlhasher"

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
	want        want
}

type mockShortNameGenerator struct{}

func (rg mockShortNameGenerator) Generate() string {
	return `mockStr`
}

func TestAddUrl(t *testing.T) {

	urlhasher.ShortNameGenerator = new(mockShortNameGenerator)

	tests := []testData{
		{
			name:        "generate short url success",
			requestBody: []byte(`https://practicum.yandex.ru/test`),
			URL:         `/`,
			method:      http.MethodPost,
			headers:     map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			want: want{
				code:     http.StatusCreated,
				response: config.GetOptions().BaseHost + `/` + urlhasher.ShortNameGenerator.Generate(),
				headers:  map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			},
		},
		{
			name:        "generate short url success with gzip compression",
			requestBody: []byte(`https://practicum.yandex.ru/test`),
			URL:         `/`,
			method:      http.MethodPost,
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding:  httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:     httphandler.HeaderContentTypeXGzip,
				httphandler.HeaderContentEncoding: httphandler.HeaderContentEncodingGzip,
			},
			want: want{
				code:     http.StatusCreated,
				response: config.GetOptions().BaseHost + `/` + urlhasher.ShortNameGenerator.Generate(),
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
			httphandler.AddURL(resp, req)

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

func TestGetUrl(t *testing.T) {

	urlhasher.ShortNameGenerator = new(mockShortNameGenerator)

	repository.GetRepository().Add(`https://practicum.yandex.ru/test`)

	tests := []testData{
		{
			name:    "get url by short name success",
			URL:     `/` + urlhasher.ShortNameGenerator.Generate(),
			method:  http.MethodGet,
			headers: map[string]string{httphandler.HeaderContentType: httphandler.HeaderContentTypeTextPlain},
			want: want{
				code:    http.StatusTemporaryRedirect,
				headers: map[string]string{httphandler.HeaderLocation: `https://practicum.yandex.ru/test`},
			},
		},
		{
			name:   "get url by short name success with gzip compression",
			URL:    `/` + urlhasher.ShortNameGenerator.Generate(),
			method: http.MethodGet,
			headers: map[string]string{
				httphandler.HeaderAcceptEncoding: httphandler.HeaderContentEncodingGzip,
				httphandler.HeaderContentType:    httphandler.HeaderContentTypeTextPlain,
			},
			want: want{
				code:    http.StatusTemporaryRedirect,
				headers: map[string]string{httphandler.HeaderLocation: `https://practicum.yandex.ru/test`},
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.URL, nil)

			for hName, hVal := range test.headers {
				req.Header.Set(hName, hVal)
			}

			resp := httptest.NewRecorder()
			httphandler.GetURL(resp, req)

			res := resp.Result()

			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			respBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)

			assert.Equal(t, test.want.response, string(respBody))
			assert.Equal(t, test.want.headers[httphandler.HeaderContentType], res.Header.Get(httphandler.HeaderContentType))
		})
	}
}
