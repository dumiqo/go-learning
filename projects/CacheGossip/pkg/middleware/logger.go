package middleware

import (
	"CacheGossip/pkg/logger"
	"net/http"
	"time"
)

func Logger(l *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Создаем обертку для перехвата статуса ответа
			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Вызываем следующий обработчик
			next.ServeHTTP(rw, r)

			// Теперь у нас есть статус и время
			duration := time.Since(start)

			l.Info("%s %s - %d %v",
				r.Method,
				r.URL.Path,
				rw.statusCode,
				duration,
			)
		})
	}
}
