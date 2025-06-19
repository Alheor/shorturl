// Package compress - сервис работы со сжатым HTTP трафиком.
//
// # Описание
//
// Используется как конвейер при обработке HTTP запросов.
package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/Alheor/shorturl/internal/httphandler"
	"github.com/Alheor/shorturl/internal/logger"
)

var _ io.Writer = (*gzipWriter)(nil)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write реализация интерфейса Writer
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHTTPHandler Gzip обработчик сжатых запросов.
func GzipHTTPHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if !strings.Contains(req.Header.Get(httphandler.HeaderAcceptEncoding), httphandler.HeaderContentEncodingGzip) {
			f(resp, req)
			return
		}

		ct := req.Header.Get(httphandler.HeaderContentType)
		if ct != httphandler.HeaderContentTypeJSON && ct != httphandler.HeaderContentTypeTextHTML && ct != httphandler.HeaderContentTypeXGzip {
			f(resp, req)
			return
		}

		if ct == httphandler.HeaderContentTypeXGzip {
			var data []byte

			data, err := io.ReadAll(req.Body)
			if err != nil {
				f(resp, req)
				return
			}

			data, err = GzipDecompress(data)
			if err != nil {
				logger.Error(`gzip decompress error:`, err)
				f(resp, req)
				return
			}

			req.Body = io.NopCloser(bytes.NewReader(data))
		}

		gz, err := gzip.NewWriterLevel(resp, gzip.BestSpeed)
		if err != nil {
			f(resp, req)
			logger.Error(`gzip error:`, err)
			return
		}

		defer gz.Close()

		resp.Header().Set(httphandler.HeaderContentEncoding, httphandler.HeaderContentEncodingGzip)

		f(gzipWriter{resp, gz}, req)
	}
}

// Compress сжатие данных.
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w, err := gzip.NewWriterLevel(&b, flate.BestSpeed)
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

// GzipDecompress разжатие данных.
func GzipDecompress(data []byte) ([]byte, error) {

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
