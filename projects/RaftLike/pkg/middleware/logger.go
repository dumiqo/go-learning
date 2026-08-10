package middleware

import (
	"RaftLike/pkg/logger"
	"net/http"
	"time"
)

func Logger(l *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			duration := time.Since(start)

			l.Info("%s %s - %v",
				r.Method,
				r.URL.Path,
				duration,
			)
		})
	}
}
