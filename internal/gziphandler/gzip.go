package gziphandler

import (
	"bytes"
	"compress/flate"
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
		if ct != `application/json` && ct != `text/html; charset=utf-8` {
			f(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			f(w, r)
			return
		}
		defer gz.Close()

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			f(w, r)
			return
		}

		data, err := Decompress(reqBody)
		if err != nil {
			f(w, r)
			return
		}

		r.Body = io.NopCloser(strings.NewReader(string(data)))

		w.Header().Set("Content-Encoding", "gzip")

		f(gzipWriter{w, gz}, r)
	}
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w, err := gzip.NewWriterLevel(&b, flate.BestCompression)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {

	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var b bytes.Buffer

	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
