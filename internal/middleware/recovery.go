package middleware

import (
	"log/slog"
	"net/http"
	response "rating/internal/transport/http"
)

func RecoveryMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("panic recovered", slog.Any("error", err))
					response.ResponseErr(log, w, http.StatusInternalServerError, "server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
