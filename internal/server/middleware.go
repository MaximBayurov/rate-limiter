package server

import (
	"net/http"
	"time"

	"github.com/MaximBayurov/rate-limiter/internal/logger"
)

func loggingMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				lrw := newLoggingResponseWriter(w)
				next.ServeHTTP(lrw, r)

				duration := time.Since(start)

				ip := r.RemoteAddr
				if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
					ip = forwarded
				}

				method := r.Method
				path := r.URL.Path
				protocol := r.Proto
				statusCode := lrw.statusCode
				latency := duration.Milliseconds()

				userAgent := r.Header.Get("User-Agent")
				if userAgent == "" {
					userAgent = "-"
				}

				(*logger).Info("request handling completed",
					"ip", ip,
					"method", method,
					"path", path,
					"protocol", protocol,
					"statusCode", statusCode,
					"latency", latency,
					"userAgent", userAgent,
				)
			},
		)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
