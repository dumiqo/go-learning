package api

import (
	"CacheGossip/internal/cache"
	"CacheGossip/pkg/logger"
	"CacheGossip/pkg/response"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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
func (c *ClientApi) Post(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, c.Name, http.StatusBadRequest,
			fmt.Errorf("invalid request body"))
		return
	}
	err := c.cache.Set(req.Key, req.Value)
	if err != nil {
		response.SendError(w, c.Name, 400, err)
	} else {
		response.SendOK(w, c.Name, "Success")
	}
}
func (c *ClientApi) Get(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	v, exist := c.cache.Get(key)
	if !exist {
		response.SendOK(w, c.Name, "not found")
	} else {
		response.SendOK(w, c.Name, v)
	}
}
func (c *ClientApi) Delete(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	err := c.cache.Remove(key)
	if err != nil {
		response.SendError(w, c.Name, 400, err)
	} else {
		response.SendOK(w, c.Name, "Success")
	}
}
