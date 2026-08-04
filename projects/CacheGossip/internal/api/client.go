package api

import (
	"CacheGossip/internal/cache"
	"CacheGossip/pkg/logger"
	"CacheGossip/pkg/response"
	"fmt"
	"net/http"
)

type ClientApi struct {
	Name   string
	cache  *cache.Cache
	logger *logger.Logger
}

func NewClientApi(name string, cache *cache.Cache, logger *logger.Logger) (*ClientApi, error) {
	if name == "" {
		return nil, fmt.Errorf("Name cannot be empty")
	}
	if cache == nil {
		return nil, fmt.Errorf("Cache cannot be nil")
	}
	return &ClientApi{name, cache, logger}, nil
}

func (c *ClientApi) Health(w http.ResponseWriter, r *http.Request) {
	response.SendOK(w, c.Name, map[string]string{"status": "ok"})
}
