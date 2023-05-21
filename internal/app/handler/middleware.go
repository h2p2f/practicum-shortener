package handler

import (
	"net/http"
	"strings"
)

// GzipHanle is middleware for gzip
func GzipHanle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentEncoding := r.Header.Get("Content-Encoding")
		supportGzip := strings.Contains(acceptEncoding, "gzip")
		sendGzip := strings.Contains(contentEncoding, "gzip")
		// this section for ordinary request
		if !supportGzip && !sendGzip {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
			return
		}
		// this section for request with accept-encoding: gzip
		if supportGzip && !sendGzip {
			originWriter := w
			compressedWriter := NewCompressWriter(w)

			originWriter = compressedWriter
			originWriter.Header().Set("Content-Encoding", "gzip")
			defer compressedWriter.Close()

			next.ServeHTTP(originWriter, r)
		}
		// this section for request with content-encoding: gzip
		if sendGzip {
			originWriter := w
			compressedWriter := NewCompressWriter(w)
			originWriter = compressedWriter
			originWriter.Header().Set("Content-Encoding", "gzip")
			defer compressedWriter.Close()

			compressedReader, err := NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			r.Body = compressedReader
			defer compressedReader.Close()

			next.ServeHTTP(originWriter, r)

		}
	})

}
