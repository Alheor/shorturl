package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/gziphandler"
	"github.com/Alheor/shorturl/internal/repository"
	"github.com/Alheor/shorturl/internal/userauth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type want struct {
	code         int
	responseBody string
	headerName   string
	headerValue  string
	cookieName   string
	cookieValue  string
}

type test struct {
	name           string
	requestURL     string
	requestBody    []byte
	cookie         *http.Cookie
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

var user = &userauth.User{ID: `5e31ae53-a6fc-43bd-8e7c-5ca06e1b413e`}

func TestAddURLSuccess(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

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
				cookieName:   userauth.CookiesName,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeTextPlainValue,
			},
		},
	}

	runTests(t, tests)
}

func TestGetURLSuccess(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

	shortName := randomShortName.Generate()
	_ = shortNameRepository.Add(ctx, user, shortName, targetURL)

	tests := []test{
		{
			name:           `positive test: call GET`,
			requestURL:     `/` + shortName,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			cookie:         prepareCookie(),
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

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

	tests := []test{
		{
			name:           `negative test: send POST without body`,
			requestURL:     `/`,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			method:         http.MethodPost,
			want: want{
				code:         http.StatusBadRequest,
				responseBody: ErrEmptyURL + "\n",
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeTextPlainValue,
			},
		},
	}

	runTests(t, tests)

	//test if exists url
	_ = shortNameRepository.Add(ctx, user, `newName`, targetURL)

	tests = []test{
		{
			name:           `negative test: send POST with existed url`,
			requestURL:     `/`,
			requestBody:    []byte(targetURL),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeTextPlainValue},
			cookie:         prepareCookie(),
			method:         http.MethodPost,
			want: want{
				code:         http.StatusConflict,
				responseBody: strings.TrimRight(config.Options.BaseHost, `/`) + `/` + randomShortName.Generate(),
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeTextPlainValue,
			},
		},
	}

	runTests(t, tests)
}

func TestGetURLError(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

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
				cookieName:   userauth.CookiesName,
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

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

	shortName := randomShortName.Generate()
	shortNameRepository.Remove(ctx, user, shortName)

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
				cookieName:   userauth.CookiesName,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		},
	}

	runTests(t, tests)
}

func TestAiShortenError(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

	shortName := randomShortName.Generate()
	shortNameRepository.Remove(ctx, user, shortName)
	_ = shortNameRepository.Add(ctx, user, shortName, targetURL)

	tests := []test{
		{
			name:           `negative api test: with existed url`,
			requestURL:     `/api/shorten`,
			requestBody:    []byte(`{"url":"` + targetURL + `"}`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			cookie:         prepareCookie(),
			method:         http.MethodPost,
			want: want{
				code:         http.StatusConflict,
				responseBody: `{"result":"` + strings.TrimRight(config.Options.BaseHost, `/`) + `/` + randomShortName.Generate() + `"}`,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		},
	}

	runTests(t, tests)
}

func TestGzip(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

	shortName := randomShortName.Generate()
	shortNameRepository.Remove(ctx, user, shortName)

	tests := []test{
		{
			name:       `test: gzip is enable on /`,
			requestURL: `/`,
			requestHeaders: map[string]string{
				HeaderAcceptEncodingName: HeaderAcceptEncodingValue,
				HeaderContentTypeName:    HeaderContentTypeTextHTMLValue,
			},
			cookie:      prepareCookie(),
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
			cookie:      prepareCookie(),
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

	shortNameRepository.Remove(ctx, user, shortName)

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

	shortNameRepository.Remove(ctx, user, shortName)

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

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

	tests := []test{
		{
			name:       `positive test: db connection /ping`,
			requestURL: `/ping`,
			method:     http.MethodGet,
			want: want{
				code:       http.StatusOK,
				cookieName: userauth.CookiesName,
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
				code:       http.StatusInternalServerError,
				cookieName: userauth.CookiesName,
			},
		},
	}

	runTests(t, tests)
}

func TestAiShortenBatchSuccess(t *testing.T) {

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

	shortName := randomShortName.Generate()
	shortNameRepository.Remove(ctx, user, shortName)

	requestBody := []byte(`[{"correlation_id":"1","original_url":"` + targetURL + `"}]`)
	responseBody := `[{"correlation_id":"1","short_url":"` + strings.TrimRight(config.Options.BaseHost, `/`) + `/` + randomShortName.Generate() + `"}]`

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
				cookieName:   userauth.CookiesName,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		},
	}

	runTests(t, tests)
}

func TestAiShortenBatchError(t *testing.T) {

	//config.Options.DatabaseDsn = "host=localhost port=5432 user=app password=pass dbname=shortener_test sslmode=disable"

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

	shortName := randomShortName.Generate()
	shortNameRepository.Remove(ctx, user, shortName)
	_ = shortNameRepository.AddBatch(ctx, user, []repository.BatchEl{{CorrelationID: "1", OriginalURL: targetURL, ShortURL: shortName}})

	tests := []test{
		{
			name:           `negative api batch test: with existed url`,
			requestURL:     `/api/shorten/batch`,
			requestBody:    []byte(`[{"correlation_id":"1","original_url":"` + targetURL + `"}]`),
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			cookie:         prepareCookie(),
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
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
				cookieName:   userauth.CookiesName,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		},
	}

	runTests(t, tests)
}

func TestAddAndGetURLForUserSuccess(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

	shortName := randomShortName.Generate()
	_ = shortNameRepository.Add(ctx, user, shortName, targetURL)
	short := strings.TrimRight(config.Options.BaseHost, `/`) + `/` + randomShortName.Generate()

	result := `[{"original_url":"https://practicum.yandex.ru/","short_url":"` + short + `"}]`

	tests := []test{
		{
			name:           `positive test: send GET`,
			requestURL:     `/api/user/urls`,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			cookie:         prepareCookie(),
			method:         http.MethodGet,
			want: want{
				code:         http.StatusOK,
				responseBody: result,
				headerName:   HeaderContentTypeName,
				headerValue:  HeaderContentTypeJSONValue,
			},
		},
	}

	runTests(t, tests)
}

func TestAddAndGetURLForUserEmptyListSuccess(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)
	shortName := randomShortName.Generate()
	shortNameRepository.Remove(ctx, user, shortName)

	tests := []test{
		{
			name:           `positive test: send GET`,
			requestURL:     `/api/user/urls`,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			cookie:         prepareCookie(),
			method:         http.MethodGet,
			want: want{
				//Все ради тестов
				//code:        http.StatusNoContent,
				code:        http.StatusUnauthorized,
				headerName:  HeaderContentTypeName,
				headerValue: HeaderContentTypeJSONValue,
			},
		},
	}

	runTests(t, tests)
}

func TestAddAndGetURLForUserUnknownUserError(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	shortNameRepository = repository.Init(ctx)

	shortName := randomShortName.Generate()
	_ = shortNameRepository.Add(ctx, user, shortName, targetURL)

	user = &userauth.User{ID: `5e3`}
	cookiesValue := string(userauth.GetSignature(user.ID))

	unknownUserCookie := &http.Cookie{
		Name:  userauth.CookiesName,
		Value: base64.StdEncoding.EncodeToString([]byte(cookiesValue)),
	}

	tests := []test{
		{
			name:           `negative test: send GET`,
			requestURL:     `/api/user/urls`,
			requestHeaders: map[string]string{HeaderContentTypeName: HeaderContentTypeJSONValue},
			cookie:         unknownUserCookie,
			method:         http.MethodGet,
			want: want{
				code: http.StatusUnauthorized,
			},
		},
	}

	runTests(t, tests)

	shortNameRepository.Remove(ctx, user, shortName)
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

			if test.cookie != nil {
				req.AddCookie(test.cookie)
			}

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

			if test.want.cookieValue != `` || test.want.cookieName != `` {
				cookieExists := false
				var cookieValue *http.Cookie
				for _, value := range resp.Cookies() {
					if value.Name == test.want.cookieName {
						cookieValue = value
						cookieExists = true
					}
				}

				if test.want.cookieValue != `` {
					assert.Equal(t, test.want.cookieValue, cookieValue.Value)
				} else {
					assert.True(t, cookieExists)
				}
			}

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

func prepareCookie() *http.Cookie {
	cookiesValue := string(userauth.GetSignature(user.ID)) + user.ID

	return &http.Cookie{
		Name:  userauth.CookiesName,
		Value: base64.StdEncoding.EncodeToString([]byte(cookiesValue)),
	}
}
