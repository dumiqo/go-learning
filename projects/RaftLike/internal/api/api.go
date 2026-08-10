package api

import (
	"RaftLike/pkg/logger"
	"RaftLike/pkg/response"
	"net/http"
)

type RaftApi struct {
	Name   string
	logger *logger.Logger
}

func NewRaftApi(name string, logger *logger.Logger) *RaftApi {
	return &RaftApi{name, logger}
}

func (g *RaftApi) Health(w http.ResponseWriter, r *http.Request) {
	response.SendOK(w, g.Name, map[string]string{"status": "ok"})
}
