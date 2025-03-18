package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/handler"
	"github.com/Alheor/shorturl/internal/urlhasher"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApiAddUrlSuccess(t *testing.T) {

	urlhasher.ShortNameGenerator = new(mockShortNameGenerator)

	tests := []testData{
		{
			name:        `API generate short url success`,
			requestBody: `{"url":"https://practicum.yandex.ru/test"}`,
			headers:     map[string]string{handler.HeaderContentTypeName: handler.HeaderContentTypeJSONValue},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusCreated,
				response: `{"result":"` + config.GetOptions().BaseHost + `/` + urlhasher.ShortNameGenerator.Generate() + `"}`,
				headers:  map[string]string{handler.HeaderContentTypeName: handler.HeaderContentTypeJSONValue},
			},
		},
		{
			name:        `API generate short with empty body error`,
			requestBody: ``,
			headers:     map[string]string{handler.HeaderContentTypeName: handler.HeaderContentTypeJSONValue},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{handler.HeaderContentTypeName: handler.HeaderContentTypeJSONValue},
			},
		},
		{
			name:        `API generate short with empty url error`,
			requestBody: `{"url":""}`,
			headers:     map[string]string{handler.HeaderContentTypeName: handler.HeaderContentTypeJSONValue},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{handler.HeaderContentTypeName: handler.HeaderContentTypeJSONValue},
			},
		},
		{
			name:        `API generate short without url field error`,
			requestBody: `{"url_test":""}`,
			headers:     map[string]string{handler.HeaderContentTypeName: handler.HeaderContentTypeJSONValue},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{handler.HeaderContentTypeName: handler.HeaderContentTypeJSONValue},
			},
		}, {
			name:        `API generate short with empty json doc error`,
			requestBody: `{}`,
			headers:     map[string]string{handler.HeaderContentTypeName: handler.HeaderContentTypeJSONValue},
			method:      http.MethodPost,
			URL:         `/api/shorten`,
			want: want{
				code:     http.StatusBadRequest,
				response: `{"error":"url required"}`,
				headers:  map[string]string{handler.HeaderContentTypeName: handler.HeaderContentTypeJSONValue},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(test.method, test.URL, bytes.NewReader([]byte(test.requestBody)))

			for hName, hVal := range test.headers {
				req.Header.Set(hName, hVal)
			}

			resp := httptest.NewRecorder()
			handler.AddShorten(resp, req)

			res := resp.Result()

			assert.Equal(t, test.want.code, res.StatusCode)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.headers[handler.HeaderContentTypeName], res.Header.Get(handler.HeaderContentTypeName))
		})
	}
}
