package utils

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var allowedContentTypes = map[string]bool{
	"application/json":    true,
	"text/html; charset:": true,
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write writes a given bytes slice to an internal buffer.
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read reads bytes from internal buffer and puts them into given slice, returning read bytes number and optional error.
func (c *gzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes internal buffer.
func (c *gzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// GzipHandle performs compression/decompression of requests and responses of given HTTP-handler.
func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		_, compressibleContentType := allowedContentTypes[r.Header.Get("Content-Type")]
		supportsGzip := compressibleContentType && !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		if supportsGzip {
			gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
			if err != nil {
				_, _ = io.WriteString(w, err.Error())
				return
			}

			ow = gzipWriter{ResponseWriter: w, Writer: gz}
			ow.Header().Set("Content-Encoding", "gzip")

			defer gz.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
