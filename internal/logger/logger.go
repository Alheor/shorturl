package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size

	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

var logger *zap.Logger

// Init Инициализация логгера
func Init(cfg *zap.Config) error {
	if logger != nil {
		return nil
	}

	var config zap.Config

	if cfg == nil {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	} else {
		config = *cfg
	}

	var err error
	logger, err = config.Build()
	if err != nil {
		return err
	}

	defer logger.Sync()

	return nil
}

// LoggingHTTPHandler логирование http запросов
func LoggingHTTPHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		start := time.Now()
		uri := req.RequestURI
		method := req.Method

		rd := &responseData{size: 0, status: 0}
		lrw := loggingResponseWriter{
			ResponseWriter: resp,
			responseData:   rd,
		}

		f(&lrw, req)

		duration := time.Since(start).String()

		logger.Info(`incoming request`,
			zap.String("url", uri),
			zap.String("method", method),
			zap.String("duration", duration),
			zap.Int("status", lrw.responseData.status),
			zap.Int("size", lrw.responseData.size),
		)
	}
}

// Info info level
func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
	defer logger.Sync()
}

// Error error level
func Error(msg string, err error) {
	if err != nil {
		logger.Error(msg + `: ` + err.Error())
	} else {
		logger.Error(msg)
	}

	defer logger.Sync()
}

// Fatal error level
func Fatal(msg string, err error) {
	if err != nil {
		logger.Error(msg + `: ` + err.Error())
	} else {
		logger.Error(msg)
	}

	logger.Sync()

	logger.Fatal(`End`)
}
