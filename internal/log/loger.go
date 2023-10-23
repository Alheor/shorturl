package log

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Logger struct {
	Log *zap.Logger
}

func Init(level string) *Logger {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	defer zl.Sync()

	lg := new(Logger)
	lg.Log = zl

	return lg
}

func (l *Logger) WithLogging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestURI := r.RequestURI
		requestMethod := r.Method
		responseData := &responseInfo{
			statusCode: 0,
			size:       0,
		}

		w1 := &loggingResponseWriter{
			ResponseWriter: w,
			ResponseInfo:   responseData,
		}

		start := time.Now()
		f(w1, r)
		duration := time.Since(start).String()

		l.Log.Info(`got incoming HTTP request`,
			zap.String("url", requestURI),
			zap.String("method", requestMethod),
			zap.String("duration", duration),
			zap.Int("status", w1.ResponseInfo.statusCode),
			zap.Int("size", w1.ResponseInfo.size),
		)
	}
}
