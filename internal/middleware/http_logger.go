package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			responseWriter := responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			start := time.Now()
			next.ServeHTTP(&responseWriter, r)
			log.Info("handler log",
				slog.String("method", r.Method),
				slog.String("url", r.URL.Path),
				slog.Duration("duration", time.Since(start)),
				slog.Int("status code", responseWriter.statusCode),
			)
		})
	}
}
