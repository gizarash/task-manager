package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriterWrapper) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
		wrappedRW := &ResponseWriterWrapper{
			ResponseWriter: w,
			StatusCode: http.StatusOK,
		}
    next.ServeHTTP(wrappedRW, r)
		slog.Info(fmt.Sprintf("%s %s %d %dms", r.Method, r.URL.Path, wrappedRW.StatusCode, time.Since(start).Milliseconds()))
  })
}