// pkg/response/response.go
package response

import (
	"encoding/json"
	"net/http"
	"time"
)

type APIResponse struct {
	Status    string      `json:"status"`
	Timestamp string      `json:"timestamp"`
	Service   string      `json:"service"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// SendJSON - отправляет стандартизированный JSON ответ
func SendJSON(w http.ResponseWriter, serviceName string, statusCode int, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   serviceName,
	}

	if err != nil {
		response.Status = "error"
		response.Error = err.Error()
	} else {
		response.Status = "ok"
		response.Data = data
	}

	// Игнорируем ошибку при кодировании
	_ = json.NewEncoder(w).Encode(response)
}

// SendError - отправляет ошибку
func SendError(w http.ResponseWriter, serviceName string, statusCode int, err error) {
	SendJSON(w, serviceName, statusCode, nil, err)
}

// SendOK - отправляет успешный ответ
func SendOK(w http.ResponseWriter, serviceName string, data interface{}) {
	SendJSON(w, serviceName, http.StatusOK, data, nil)
}
