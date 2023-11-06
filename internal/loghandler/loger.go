package loghandler

import (
	"github.com/Alheor/shorturl/internal/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Logger struct {
	Log *zap.Logger
}

func Init() *Logger {
	lvl, err := zap.ParseAtomicLevel(config.Options.LogLevel)
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

		lw := &loggingResponseWriter{
			ResponseWriter: w,
			ResponseInfo:   responseData,
		}

		start := time.Now()
		f(lw, r)
		duration := time.Since(start).String()

		encodingType := lw.Header().Get(`Content-Encoding`)
		if encodingType != `` {
			encodingType = ` with encoding: ` + encodingType
		}

		l.Log.Info(`got incoming HTTP request`+encodingType,
			zap.String("url", requestURI),
			zap.String("method", requestMethod),
			zap.String("duration", duration),
			zap.Int("status", lw.ResponseInfo.statusCode),
			zap.Int("size", lw.ResponseInfo.size),
		)
	}
}
