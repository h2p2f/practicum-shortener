package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

var Log *zap.Logger = zap.NewNop()

type responseData struct {
	status int
	size   int
}
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
func InitLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	Log, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{}
		loggedw := loggingResponseWriter{w, responseData}
		h.ServeHTTP(&loggedw, r)
		//Log.Info("Request", zap.String("url", r.URL.String()), zap.String("method", r.Method), zap.Duration("duration", time.Since(start)))
		//Log.Info("Response", zap.Int("status", responseData.status), zap.Int("size", responseData.size))
		Log.Sugar().Infof("Request  - method: %s, url: %s, duration: %s", r.Method, r.URL.String(), time.Since(start))
		Log.Sugar().Infof("Response - status: %d, size: %d", responseData.status, responseData.size)
	}
	return http.HandlerFunc(logFn)
}
