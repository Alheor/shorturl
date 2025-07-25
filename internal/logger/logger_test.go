package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type logResult struct {
	URL      string `json:"url"`
	Method   string `json:"method"`
	Duration string `json:"duration"`
	Status   int    `json:"status"`
	Size     int    `json:"size"`
}

func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }

type test struct {
	name         string
	requestURL   string
	method       string
	responseCode int
	responseSize int
}

func TestLogging(t *testing.T) {
	tests := []test{
		{
			name:         `test 1`,
			requestURL:   `/test-request`,
			method:       http.MethodGet,
			responseCode: http.StatusOK,
			responseSize: 3,
		},
		{
			name:         `test 2`,
			requestURL:   `/`,
			method:       http.MethodPost,
			responseCode: http.StatusBadRequest,
			responseSize: 12,
		},
	}

	runTests(t, tests)
}

func TestInitWithDefaultCfg(t *testing.T) {

	err := Init(nil)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestInitWithOtherCfg(t *testing.T) {

	cfg := zap.NewProductionConfig()

	err := Init(&cfg)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func runTests(t *testing.T, tests []test) {
	sink := &MemorySink{new(bytes.Buffer)}

	err := zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
		return sink, nil
	})
	require.NoError(t, err)

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"memory://"}
	err = Init(&cfg)
	require.NoError(t, err)

	ts := httptest.NewServer(getRoutes())
	defer ts.Close()

	client := ts.Client()

	for _, test := range tests {
		sink.Reset()

		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, ts.URL+test.requestURL, nil)
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)

			_, err = io.ReadAll(resp.Body)
			defer resp.Body.Close()
			require.NoError(t, err)

			result := new(logResult)
			err = json.Unmarshal(sink.Bytes(), result)
			require.NoError(t, err)

			assert.Equal(t, test.requestURL, result.URL)
			assert.Equal(t, test.method, result.Method)
			assert.Equal(t, test.responseCode, result.Status)
			assert.Equal(t, test.responseSize, result.Size)
			assert.NotEmpty(t, result.Duration)
		})
	}
}

func getRoutes() chi.Router {
	r := chi.NewRouter()

	testGetHandler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `ok`, http.StatusOK)
	}

	testPostHandler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `bad request`, http.StatusBadRequest)
	}

	r.Get(`/*`, LoggingHTTPHandler(testGetHandler))
	r.Post(`/`, LoggingHTTPHandler(testPostHandler))

	return r
}
