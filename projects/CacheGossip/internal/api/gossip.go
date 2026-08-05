package api

import (
	"CacheGossip/internal/gossip"
	"CacheGossip/pkg/logger"
	"CacheGossip/pkg/models"
	"CacheGossip/pkg/response"
	"encoding/json"
	"fmt"
	"net/http"
)

type GossipApi struct {
	Name   string
	gossip *gossip.Gossip
	logger *logger.Logger
}

func NewGossipApi(name string, membersip *gossip.Gossip, logger *logger.Logger) *GossipApi {
	return &GossipApi{name, membersip, logger}
}

func (g *GossipApi) Health(w http.ResponseWriter, r *http.Request) {
	response.SendOK(w, g.Name, map[string]string{"status": "ok"})
}

func (g *GossipApi) Gossip(w http.ResponseWriter, r *http.Request) {
	var req models.GossipMessage

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, g.Name, http.StatusBadRequest, fmt.Errorf("invalid request body %s", err))
		return
	}
	if err := req.Validate(); err != nil {
		response.SendError(w, g.Name, http.StatusBadRequest, err)
		return
	}
	g.logger.Info("!!!!!!!!!!!!! gossip !!!!!!!!!!!!!")
}
