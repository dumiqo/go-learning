package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Type int

const (
	Membership Type = iota
)

type GossipMessage struct {
	Sender string    `json:"sender"`
	Time   time.Time `json:"time"`
	UUID   uuid.UUID `json:"uuid"`
	Type   Type      `json:"type"`

	Nodes []NodeInfo `json:"nodes,omitempty"`
}

type NodeInfo struct {
	Name     string    `json:"name"`
	Status   Status    `json:"status"`
	Address  string    `json:"address"`
	LastSeen time.Time `json:"last_seen"`
}

func (r *GossipMessage) Validate() error {
	if r.Sender == "" {
		return fmt.Errorf("Sender is required")
	}
	if r.UUID == uuid.Nil {
		return fmt.Errorf("uuid is required")
	}
	if time.Now().Before(r.Time) {
		return fmt.Errorf("gossip from future")
	}
	return nil
}
func (r *GossipMessage) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
func (r *GossipMessage) FromJSON(data []byte) error {
	return json.Unmarshal(data, r)
}
func NewMembershipGossip(name string, nodes []NodeInfo) GossipMessage {
	return GossipMessage{
		Sender: name,
		UUID:   uuid.New(),
		Time:   time.Now().UTC(),
		Type:   Membership,

		Nodes: nodes,
	}
}
