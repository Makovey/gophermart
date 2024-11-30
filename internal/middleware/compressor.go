package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"
	"strings"
)

var allowedCompress = []string{
	"application/json",
	"text/html",
}

type Compressor struct{}

func (c Compressor) Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		if slices.Contains(allowedCompress, acceptEncoding) {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		if strings.Contains(contentEncoding, "gzip") {
			gr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			r.Body = gr
			gr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}

func NewCompressor() Compressor {
	return Compressor{}
}

type compressWriter struct {
	w  http.ResponseWriter
	gw *gzip.Writer
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.gw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}

	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Close() error {
	return c.gw.Close()
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{w: w, gw: gzip.NewWriter(w)}
}

type compressReader struct {
	r  io.ReadCloser
	gr *gzip.Reader
}

func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.gr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.gr.Close()
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressReader{r: r, gr: gr}, nil
}
