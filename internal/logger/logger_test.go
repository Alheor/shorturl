package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Alheor/shorturl/internal/config"
	"github.com/Alheor/shorturl/internal/handler"
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

type MemorySink struct {
	*bytes.Buffer
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
			responseCode: http.StatusBadRequest,
			responseSize: 19,
		},
		{
			name:         `test 2`,
			requestURL:   `/`,
			method:       http.MethodPost,
			responseCode: http.StatusBadRequest,
			responseSize: 13,
		},
	}

	runTests(t, tests)
}

func runTests(t *testing.T, tests []test) {
	sink := &MemorySink{new(bytes.Buffer)}

	err := zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
		return sink, nil
	})
	require.NoError(t, err)
	config.Load()

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"memory://"}
	Init(&cfg)

	ts := httptest.NewServer(getRoutes())
	defer ts.Close()

	client := ts.Client()

	for _, test := range tests {
		sink.Reset()

		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, ts.URL+test.requestURL, nil)
			require.NoError(t, err)

			_, err = client.Do(req)
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

	r.Get(`/*`, LoggingHTTPHandler(handler.GetURL))
	r.Post(`/`, LoggingHTTPHandler(handler.AddURL))

	return r
}
