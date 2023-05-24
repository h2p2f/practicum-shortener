package handler

import (
	"compress/gzip"
	"io"
	"net/http"
)

// CompressWriter is implementation of http.ResponseWriter
type CompressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

// CompressReader is implementation of io.ReadCloser
type CompressReader struct {
	r io.ReadCloser
	z *gzip.Reader
}

// NewCompressWriter is constructor for CompressWriter
func NewCompressWriter(w http.ResponseWriter) *CompressWriter {
	zw := gzip.NewWriter(w)
	return &CompressWriter{w, zw}
}

// Header is implementation of http.ResponseWriter.Header
func (cw *CompressWriter) Header() http.Header {
	return cw.w.Header()
}

// Write is implementation of http.ResponseWriter.Write
func (cw *CompressWriter) Write(b []byte) (int, error) {
	return cw.zw.Write(b)
}

// WriteHeader is implementation of http.ResponseWriter.WriteHeader
func (cw *CompressWriter) WriteHeader(statusCode int) {
	cw.w.WriteHeader(statusCode)
	if statusCode > 199 && statusCode < 300 {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}
}

// Close is closes gzip.Writer
func (cw *CompressWriter) Close() error {
	return cw.zw.Close()
}

// NewCompressReader is constructor for CompressReader
func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	z, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &CompressReader{r, z}, nil
}

// Read is implementation of io.ReadCloser.Read
func (cr *CompressReader) Read(b []byte) (int, error) {
	return cr.z.Read(b)
}

// Close is implementation of io.ReadCloser.Close
func (cr *CompressReader) Close() error {
	if err := cr.z.Close(); err != nil {
		return err
	}
	return cr.r.Close()
}
