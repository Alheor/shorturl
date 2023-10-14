package main

import (
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/repository"
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
	requestURL  string
	requestBody io.Reader
	method      string
	want        want
}

const targetURL = `https://practicum.yandex.ru/`

type mockShortName struct{}

func (rg mockShortName) Generate() string {
	return `mockString`
}

func init() {
	randomShortName = new(mockShortName)
}

func TestAddURLSuccess(t *testing.T) {

	tests := []test{
		{
			name:        `positive test: send POST`,
			requestURL:  `/`,
			requestBody: strings.NewReader(targetURL),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusCreated,
				responseBody: strings.TrimRight(config.Options.BaseHost, `/`) + `/` + randomShortName.Generate(),
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeValue,
			},
		},
	}

	runTests(t, tests)
}

func TestGetURLSuccess(t *testing.T) {
	shortName := randomShortName.Generate()
	_ = shortNameRepository.Add(shortName, targetURL)

	tests := []test{
		{
			name:       `positive test: call GET`,
			requestURL: `/` + shortName,
			method:     http.MethodGet,
			want: want{
				code:        http.StatusTemporaryRedirect,
				headerName:  HeaderLocation,
				headerValue: targetURL,
			},
		},
	}

	runTests(t, tests)
}

func TestAddURLError(t *testing.T) {

	tests := []test{
		{
			name:       `negative test: send POST without body`,
			requestURL: `/`,
			method:     http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: ErrorEmptyRequestBody + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeValue,
			},
		}, {
			name:        `negative test: send POST empty body 1`,
			requestURL:  `/`,
			requestBody: strings.NewReader(``),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: ErrorEmptyRequestBody + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeValue,
			},
		}, {
			name:        `negative test: send POST empty body 2`,
			requestURL:  `/`,
			requestBody: strings.NewReader(` `),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: ErrorEmptyRequestBody + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeValue,
			},
		}, {
			name:        `negative test: send POST invalid url`,
			requestURL:  `/`,
			requestBody: strings.NewReader(`invalid url`),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: ErrorInvalidURL + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeValue,
			},
		},
	}

	runTests(t, tests)

	//test if exists by url
	_ = shortNameRepository.Add(`otherShortName`, targetURL)

	tests = []test{
		{
			name:        `negative test: send POST with existed url`,
			requestURL:  `/`,
			requestBody: strings.NewReader(targetURL),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: repository.ErrorValueAlreadyExist + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeValue,
			},
		},
	}

	runTests(t, tests)

	//test if exists by short name
	_ = shortNameRepository.Add(randomShortName.Generate(), targetURL)

	tests = []test{
		{
			name:        `negative test: send POST with existed short name`,
			requestURL:  `/`,
			requestBody: strings.NewReader(targetURL),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: repository.ErrorValueAlreadyExist + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeValue,
			},
		},
	}

	runTests(t, tests)
}

func TestGetURLError(t *testing.T) {

	tests := []test{
		{
			name:       `negative test: method GET not allowed`,
			requestURL: `/`,
			method:     http.MethodGet,
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		}, {
			name:       `negative test: call GET unknown identifier`,
			requestURL: `/unknown-identifier`,
			method:     http.MethodGet,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: repository.ErrorIdNotFound + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeValue,
			},
		}, {
			name:       `negative test: method POST not allowed`,
			requestURL: `/unknown-identifier`,
			method:     http.MethodPost,
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}

	runTests(t, tests)
}

func runTests(t *testing.T, tests []test) {

	ts := httptest.NewServer(getRouter())
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			req, err := http.NewRequest(test.method, ts.URL+test.requestURL, test.requestBody)
			require.NoError(t, err)

			client := ts.Client()
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.headerValue, resp.Header.Get(test.want.headerName))

			//проверка тела ответа
			resBody, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.responseBody, string(resBody))
		})
	}
}
