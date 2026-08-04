package api

import (
	"CacheGossip/internal/cache"
	"CacheGossip/pgk/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ClientApi struct {
	Name   *string
	cache  *cache.Cache
	logger *logger.Logger
}

func NewClientApi(name *string, cache *cache.Cache, logger *logger.Logger) (*ClientApi, error) {
	if name == nil {
		return nil, fmt.Errorf("Name cannot be empty")
	}
	if cache == nil {
		return nil, fmt.Errorf("Cache cannot be nil")
	}
	return &ClientApi{name, cache, logger}, nil
}

func (c *ClientApi) Health(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Health check")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем заголовок
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Отправляем JSON ответ
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "CacheGossip",
	}

	json.NewEncoder(w).Encode(response)
}
