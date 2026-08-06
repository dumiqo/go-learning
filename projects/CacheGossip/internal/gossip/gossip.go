package gossip

import (
	"CacheGossip/internal/membership"
	"CacheGossip/pkg/logger"
	"CacheGossip/pkg/models"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

type Gossip struct {
	nodeName, address string
	membership        *membership.Membership
	logger            *logger.Logger
}
type GossipStatus struct {
	Memberships *membership.MembershipStatus
}

func NewGossip(nodeName, nodeAddress string, membership *membership.Membership, logger *logger.Logger) *Gossip {
	return &Gossip{nodeName, nodeAddress, membership, logger}
}

func (g *Gossip) Start(ctx context.Context) {
	go g.membership.Start(ctx)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	g.logger.Info("Start sending gossip")
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

func (g *Gossip) Status() *GossipStatus {
	mStatus := g.membership.Status()
	return &GossipStatus{mStatus}
}

func (g *Gossip) ProcessGossip(msg models.GossipMessage) error {
	g.logger.Info("Begin processiong gossip from: %s. uuid: %s", msg.Sender, msg.UUID)
	switch msg.Type {
	case models.Membership:
		g.ProcessMembership(msg)
	default:
		return fmt.Errorf("Unknown msg type: %s", msg.Type)
	}
	g.logger.Info("End in processiong gossip from: %s. uuid: %s", msg.Sender, msg.UUID)
	return nil
}

func (g *Gossip) ProcessMembership(msg models.GossipMessage) {
	for _, node := range msg.Nodes {
		if g.nodeName == node.Name {
			continue
		}

		if err := g.membership.Update(node); err != nil {
			g.logger.Error("Error in updating membership. member: %s. %s", node.Name, err)
		}
		g.logger.Info("Update member, %s", node.Name)
	}

	if g.membership.IsKnownMember(msg.Sender) {
		g.membership.UpdateStatus(msg.Sender, msg.Time)
		return
	}

	g.membership.AddMember(msg.Sender, msg.Address, msg.Time)
	g.logger.Info("New member %s", msg.Sender)
	return
}

func (g *Gossip) sendGossip() {
	member := g.membership.RandomMember()

	msg := models.NewMembershipGossip(g.nodeName, g.address, g.membership.GetNodeInfo())

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
		g.membership.KillNode(g.nodeName)
		g.logger.Warning("gossip response error: status %d", resp.StatusCode)
		return
	}
}
