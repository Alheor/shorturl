package loghandler

import "net/http"

type (
	responseInfo struct {
		statusCode int
		size       int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		ResponseInfo *responseInfo
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.ResponseInfo.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseInfo.statusCode = statusCode
}
