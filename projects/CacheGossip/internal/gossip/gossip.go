package gossip

import (
	"CacheGossip/internal/membership"
	"CacheGossip/pkg/logger"
	"CacheGossip/pkg/models"
	"bytes"
	"context"
	"net/http"
	"time"
)

type Gossip struct {
	nodeName   string
	membership *membership.Membership
	logger     *logger.Logger
}

func NewGossip(nodeName string, membership *membership.Membership, logger *logger.Logger) *Gossip {
	return &Gossip{nodeName, membership, logger}
}

func (g *Gossip) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	g.logger.Info("Start sending gossip")
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			g.sendGossip()
		case <-ctx.Done():
			g.logger.Info("End sending gossip")
			return
		}
	}
}

func (g *Gossip) ProcessGossip(msg models.GossipMessage) error {
	g.logger.Info("Begin processiong gossip from: %s. uuid: %s", msg.Sender, msg.UUID)
	for _, node := range msg.Nodes {
		if g.nodeName == node.Name {
			continue
		}

		if err := g.membership.Update(node); err != nil {
			g.logger.Error("Error in updating membership. member: %s. %s", node.Name, err)
		}
		g.logger.Info("Update member, %s", node.Name)
	}

	if err := g.membership.UpdateStatus(msg.Sender, msg.Time); err != nil {
		g.logger.Error("Error in updating membership. member: %s. %s", msg.Sender, err)
	}
	g.logger.Info("End in processiong gossip from: %s. uuid: %s", msg.Sender, msg.UUID)
	return nil
}

func (g *Gossip) sendGossip() {
	member := g.membership.RandomMember()

	msg := models.NewMembershipGossip(g.nodeName, g.membership.GetNodeInfo())

	json, err := msg.ToJSON()

	if err != nil {
		g.logger.Error("Error in serializing. %s", err)
		return
	}
	req, err := http.NewRequest("POST", member.Url, bytes.NewBuffer(json))
	if err != nil {
		g.logger.Error("Error creating request to %s. %s", member.Name, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		g.logger.Error("gossip request failed: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		g.logger.Warning("gossip response error: status %d", resp.StatusCode)
		return
	}
}
