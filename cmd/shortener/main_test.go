package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/gziphandler"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
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
	name           string
	requestURL     string
	requestBody    []byte
	requestHeaders map[string]string
	method         string
	want           want
}

const targetURL = `https://practicum.yandex.ru/`

type mockShortName struct{}

func (rg mockShortName) Generate() string {
	return `mockStr`
}

func init() {
	randomShortName = new(mockShortName)
	config.Options.FileStoragePath = `` //режим без записи в файл
}

func TestAddURLSuccess(t *testing.T) {

	shortNameRepository = repository.Init()
	shortName := randomShortName.Generate()
	shortNameRepository.Remove(shortName)

	tests := []test{
		{
			name:           `positive test: send POST`,
			requestURL:     `/`,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			requestBody:    []byte(targetURL),
			method:         http.MethodPost,
			want: want{
				code:         http.StatusCreated,
				responseBody: strings.TrimRight(config.Options.BaseHost, `/`) + `/` + randomShortName.Generate(),
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeTextPlainValue,
			},
		},
	}

	runTests(t, tests)
}

func TestGetURLSuccess(t *testing.T) {

	shortNameRepository = repository.Init()

	shortName := randomShortName.Generate()
	_ = shortNameRepository.Add(shortName, targetURL)

	tests := []test{
		{
			name:           `positive test: call GET`,
			requestURL:     `/` + shortName,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodGet,
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

	shortNameRepository = repository.Init()

	tests := []test{
		{
			name:           `negative test: send POST without body`,
			requestURL:     `/`,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: ErrEmptyURL + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeTextPlainValue,
			},
		}, {
			name:           `negative test: send POST empty body 1`,
			requestURL:     `/`,
			requestBody:    []byte(``),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: ErrEmptyURL + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeTextPlainValue,
			},
		}, {
			name:           `negative test: send POST empty body 2`,
			requestURL:     `/`,
			requestBody:    []byte(` `),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: ErrEmptyURL + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeTextPlainValue,
			},
		}, {
			name:           `negative test: send POST invalid url`,
			requestURL:     `/`,
			requestBody:    []byte(`invalid url`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: ErrInvalidURL + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeTextPlainValue,
			},
		},
	}

	runTests(t, tests)

	//test if exists url
	_ = shortNameRepository.Add(`otherShortName`, targetURL)

	tests = []test{
		{
			name:           `negative test: send POST with existed url`,
			requestURL:     `/`,
			requestBody:    []byte(targetURL),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: repository.ErrValueAlreadyExist + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeTextPlainValue,
			},
		},
	}

	runTests(t, tests)
}

func TestGetURLError(t *testing.T) {

	shortNameRepository = repository.Init()

	tests := []test{
		{
			name:           `negative test: method GET not allowed`,
			requestURL:     `/`,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodGet,
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		}, {
			name:           `negative test: call GET unknown identifier`,
			requestURL:     `/unknown-identifier`,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodGet,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: repository.ErrIDNotFound + "\n",
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeTextPlainValue,
			},
		}, {
			name:           `negative test: method POST not allowed`,
			requestURL:     `/unknown-identifier`,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodPost,
			want: want{
				code: http.StatusMethodNotAllowed,
			},
		},
	}

	runTests(t, tests)
}

func TestAiShortenSuccess(t *testing.T) {

	shortNameRepository = repository.Init()

	shortName := randomShortName.Generate()
	shortNameRepository.Remove(shortName)

	tests := []test{
		{
			name:           `positive api test: send POST with valid body`,
			requestURL:     `/api/shorten`,
			requestBody:    []byte(`{"url":"` + targetURL + `"}`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusCreated,
				responseBody: `{"result":"` + strings.TrimRight(config.Options.BaseHost, `/`) + `/` + randomShortName.Generate() + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		},
	}

	runTests(t, tests)
}

func TestAiShortenError(t *testing.T) {

	shortNameRepository = repository.Init()

	shortName := randomShortName.Generate()
	shortNameRepository.Remove(shortName)
	_ = shortNameRepository.Add(shortName, targetURL)

	tests := []test{
		{
			name:           `negative api test: with existed url`,
			requestURL:     `/api/shorten`,
			requestBody:    []byte(`{"url":"` + targetURL + `"}`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + repository.ErrValueAlreadyExist + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api test: invalid url`,
			requestURL:     `/api/shorten`,
			requestBody:    []byte(`{"url":"invalid url"}`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrInvalidURL + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api test: empty url 1`,
			requestURL:     `/api/shorten`,
			requestBody:    []byte(`{"url":""}`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrEmptyURL + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api test: empty url 2`,
			requestURL:     `/api/shorten`,
			requestBody:    []byte(`{"url":" "}`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrEmptyURL + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api test: url is null`,
			requestURL:     `/api/shorten`,
			requestBody:    []byte(`{"url":null}`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrEmptyURL + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api test: invalid json`,
			requestURL:     `/api/shorten`,
			requestBody:    []byte(`{"url":null`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrOnlyJSONDataAllowed + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api test: invalid header`,
			requestURL:     `/api/shorten`,
			requestBody:    []byte(`{"url":"` + targetURL + `"}`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrOnlyJSONDataAllowed + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		},
	}

	runTests(t, tests)
}

func TestGzip(t *testing.T) {

	shortNameRepository = repository.Init()

	shortName := randomShortName.Generate()
	shortNameRepository.Remove(shortName)

	tests := []test{
		{
			name:       `test: gzip is enable on /`,
			requestURL: `/`,
			requestHeaders: map[string]string{
				HeaderAcceptEncodingName: HeaderAcceptEncodingValue,
				HeaderContentTypeName:    HeaderContentTypeTextHTMLValue,
			},
			requestBody: []byte(targetURL),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusCreated,
				responseBody: strings.TrimRight(config.Options.BaseHost, `/`) + `/` + randomShortName.Generate(),
				headerName:   HeaderContentEncodingName,
				headerValue:  HeaderAcceptEncodingValue,
			},
		}, {
			name:       `test: gzip is enable on /{id}`,
			requestURL: `/` + shortName,
			requestHeaders: map[string]string{
				HeaderAcceptEncodingName: HeaderAcceptEncodingValue,
				HeaderContentTypeName:    HeaderContentTypeTextHTMLValue,
			},
			requestBody: []byte(targetURL),
			method:      http.MethodGet,
			want: want{
				code:        http.StatusTemporaryRedirect,
				headerName:  HeaderLocation,
				headerValue: targetURL,
			},
		},
	}

	runTests(t, tests)

	shortNameRepository.Remove(shortName)

	tests = []test{
		{
			name:       `test: gzip is enable on /api/shorten`,
			requestURL: `/api/shorten`,
			requestHeaders: map[string]string{
				HeaderAcceptEncodingName: HeaderAcceptEncodingValue,
				HeaderContentTypeName:    HeaderContentTypeJSONValue,
			},
			requestBody: []byte(`{"url":"` + targetURL + `"}`),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusCreated,
				responseBody: `{"result":"` + strings.TrimRight(config.Options.BaseHost, `/`) + `/` + randomShortName.Generate() + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		},
	}

	runTests(t, tests)

	shortNameRepository.Remove(shortName)

	tests = []test{
		{
			name:       `test: gzip is enable on /api/shorten with application/x-gzip header`,
			requestURL: `/api/shorten`,
			requestHeaders: map[string]string{
				HeaderAcceptEncodingName: HeaderAcceptEncodingValue,
				HeaderContentTypeName:    HeaderContentTypeXgzipValue,
			},
			requestBody: []byte(`{"url":"` + targetURL + `"}`),
			method:      http.MethodPost,
			want: want{
				code:         http.StatusCreated,
				responseBody: `{"result":"` + strings.TrimRight(config.Options.BaseHost, `/`) + `/` + randomShortName.Generate() + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		},
	}

	runTests(t, tests)
}

func TestPingSuccess(t *testing.T) {

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	shortNameRepository = repository.Init()

	tests := []test{
		{
			name:       `positive test: db connection /ping`,
			requestURL: `/ping`,
			method:     http.MethodGet,
			want: want{
				code: http.StatusOK,
			},
		},
	}

	runTests(t, tests)
}

func TestPingError(t *testing.T) {

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	instance := new(repository.Postgres)

	conn, err := pgxpool.New(context.Background(), config.Options.DatabaseDsn)
	if err != nil {
		require.NoError(t, err)
	}
	conn.Close()

	instance.Conn = conn
	shortNameRepository = instance

	tests := []test{
		{
			name:       `negative test: db connection /ping`,
			requestURL: `/ping`,
			method:     http.MethodGet,
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	runTests(t, tests)
}

func TestAiShortenBatchSuccess(t *testing.T) {

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	shortNameRepository = repository.Init()

	shortName := randomShortName.Generate()
	shortNameRepository.Remove(shortName)

	requestBody := []byte(`[{"correlation_id":"1","original_url":"` + targetURL + `"}]`)
	responseBody := `[{"correlation_id":"1","short_url":"` + randomShortName.Generate() + `"}]`

	tests := []test{
		{
			name:           `positive api batch test: send POST with valid body`,
			requestURL:     `/api/shorten/batch`,
			requestBody:    requestBody,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusCreated,
				responseBody: responseBody,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		},
	}

	runTests(t, tests)
}

func TestAiShortenBatchError(t *testing.T) {

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	shortNameRepository = repository.Init()

	shortName := randomShortName.Generate()
	shortNameRepository.Remove(shortName)
	_ = shortNameRepository.AddBatch([]repository.BatchEl{{CorrelationID: "1", OriginalURL: targetURL, ShortURL: shortName}})

	tests := []test{
		{
			name:           `negative api batch test: with existed url`,
			requestURL:     `/api/shorten/batch`,
			requestBody:    []byte(`[{"correlation_id":"1","original_url":"` + targetURL + `"}]`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + repository.ErrValueAlreadyExist + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api batch test: invalid url`,
			requestURL:     `/api/shorten/batch`,
			requestBody:    []byte(`[{"correlation_id":"1","original_url":"invalid_url"}]`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrInvalidURL + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api batch test: empty data 1`,
			requestURL:     `/api/shorten/batch`,
			requestBody:    []byte(`[{}]`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrEmptyURL + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api batch test: empty data 2`,
			requestURL:     `/api/shorten/batch`,
			requestBody:    []byte(`[]`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrEmptyURL + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api batch test: empty data 3`,
			requestURL:     `/api/shorten/batch`,
			requestBody:    []byte(``),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrOnlyJSONDataAllowed + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api batch test: data is null`,
			requestURL:     `/api/shorten/batch`,
			requestBody:    []byte(`null`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrEmptyURL + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api test: url is null`,
			requestURL:     `/api/shorten/batch`,
			requestBody:    []byte(`[{"correlation_id":"1","original_url":null}]`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrEmptyURL + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api test: invalid json`,
			requestURL:     `/api/shorten/batch`,
			requestBody:    []byte(`{"url":null`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrOnlyJSONDataAllowed + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		}, {
			name:           `negative api test: invalid header`,
			requestURL:     `/api/shorten/batch`,
			requestBody:    []byte(`{"url":"` + targetURL + `"}`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: `{"error":"` + ErrOnlyJSONDataAllowed + `"}`,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
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

			var err error
			if test.requestHeaders[HeaderContentTypeName] == HeaderContentTypeXgzipValue {
				test.requestBody, err = gziphandler.Compress(test.requestBody)

				require.NoError(t, err)
			}

			req, err := http.NewRequest(test.method, ts.URL+test.requestURL, bytes.NewReader(test.requestBody))
			require.NoError(t, err)

			for hName, hVal := range test.requestHeaders {
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

			respHeader := resp.Header.Get(test.want.headerName)
			assert.Equal(t, test.want.headerValue, respHeader)

			var respBody io.ReadCloser

			if resp.Header.Get(HeaderContentEncodingName) == HeaderAcceptEncodingValue {
				respBody, err = gzip.NewReader(resp.Body)
				require.NoError(t, err)
			} else {
				respBody = resp.Body
			}

			//проверка тела ответа
			resBody, err := io.ReadAll(respBody)
			defer resp.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, test.want.responseBody, string(resBody))
		})
	}
}
