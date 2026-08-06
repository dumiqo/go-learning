package middleware

import (
	"bytes"
	"net/http"
)

// responseWriter - обертка для перехвата ответа
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	// Сохраняем тело ответа
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}
