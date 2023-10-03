package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type want struct {
	code         int
	responseBody string
	headerName   string
	headerValue  string
}

type test struct {
	name        string
	requestUrl  string
	requestBody io.Reader
	method      string
	want        want
}

const targetUrl = `https://practicum.yandex.ru/`

func TestAddUrlSuccess(t *testing.T) {

	tests := []test{
		{
			name:        `positive test send POST`,
			requestUrl:  `/`,
			requestBody: strings.NewReader(targetUrl),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusCreated,
				responseBody: `http://` + addr + `/` + shortName,
				headerName:   `Content-Type`,
				headerValue:  `text/plain; charset=utf-8`,
			},
		},
	}

	runTests(t, tests)
}

func TestGetUrlSuccess(t *testing.T) {

	urlMap[shortName] = targetUrl
	tests := []test{
		{
			name:       `positive test #2: call GET`,
			requestUrl: `/` + shortName,
			method:     http.MethodGet,
			want: want{
				code:        http.StatusTemporaryRedirect,
				headerName:  `Location`,
				headerValue: targetUrl,
			},
		},
	}

	runTests(t, tests)
}

func TestAddUrlError(t *testing.T) {

	tests := []test{
		{
			name:       `negative test #1: send POST without body`,
			requestUrl: `/`,
			method:     http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: "Request body is empty\n",
				headerName:   `Content-Type`,
				headerValue:  `text/plain; charset=utf-8`,
			},
		}, {
			name:        `negative test #1: send POST empty body 1`,
			requestUrl:  `/`,
			requestBody: strings.NewReader(``),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: "Request body is empty\n",
				headerName:   `Content-Type`,
				headerValue:  `text/plain; charset=utf-8`,
			},
		}, {
			name:        `negative test #1: send POST empty body 2`,
			requestUrl:  `/`,
			requestBody: strings.NewReader(` `),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: "Request body is empty\n",
				headerName:   `Content-Type`,
				headerValue:  `text/plain; charset=utf-8`,
			},
		}, {
			name:        `negative test #1: send POST invalid url`,
			requestUrl:  `/`,
			requestBody: strings.NewReader(`invalid url`),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: "Only valid url allowed\n",
				headerName:   `Content-Type`,
				headerValue:  `text/plain; charset=utf-8`,
			},
		},
	}

	runTests(t, tests)
}

func TestGetUrlError(t *testing.T) {

	tests := []test{
		{
			name:       `negative test #1: call GET empty identifier`,
			requestUrl: `/`,
			method:     http.MethodGet,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: "Identifier is empty\n",
				headerName:   `Content-Type`,
				headerValue:  `text/plain; charset=utf-8`,
			},
		}, {
			name:       `negative test #1: call GET unknown identifier`,
			requestUrl: `/unknown-identifier`,
			method:     http.MethodGet,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: "Unknown identifier\n",
				headerName:   `Content-Type`,
				headerValue:  `text/plain; charset=utf-8`,
			},
		},
	}

	runTests(t, tests)
}

func runTests(t *testing.T, tests []test) {

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.requestUrl, test.requestBody)
			w := httptest.NewRecorder()

			if test.method == http.MethodPost {
				addURL(w, request)
			} else {
				getURL(w, request)
			}

			res := w.Result()

			//првоерка кода ответа
			assert.Equal(t, test.want.code, res.StatusCode)

			//проверка заголовка
			assert.Equal(t, test.want.headerValue, res.Header.Get(test.want.headerName))

			if test.method == http.MethodGet {
				return
			}

			//проверка тела ответа
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				require.NoError(t, err)
			}(res.Body)
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.responseBody, string(resBody))
		})
	}
}
