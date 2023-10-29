package gziphandler

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func WithGzip(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get(`Accept-Encoding`), `gzip`) {
			f(w, r)
			return
		}

		ct := r.Header.Get(`Content-Type`)
		if ct != `application/json` && ct != `text/plain; charset=utf-8` {
			f(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			f(w, r)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")

		f(gzipWriter{w, gz}, r)
	}
}
