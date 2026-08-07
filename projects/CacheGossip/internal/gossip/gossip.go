package gossip

import (
	"CacheGossip/internal/cache"
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
	cache             *cache.Cache
	logger            *logger.Logger
}
type GossipStatus struct {
	Memberships *membership.MembershipStatus
}

func NewGossip(nodeName, nodeAddress string, membership *membership.Membership, cache *cache.Cache, logger *logger.Logger) *Gossip {
	return &Gossip{nodeName, nodeAddress, membership, cache, logger}
}

func (g *Gossip) Start(ctx context.Context) {
	go g.membership.Start(ctx)
	go g.startMembershipGossip(ctx)
	go g.startDataGossip(ctx)
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
	case models.Data:
		g.ProcessData(msg)
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
}

func (g *Gossip) ProcessData(msg models.GossipMessage) {
	if len(msg.Operations) <= 0 {
		g.logger.Info("Nothing to apply from: %s.", msg.Sender)
		return
	}
	for _, op := range msg.Operations {
		switch op.Type {
		case models.Set:
			g.cache.Set(op.Key, op.Value, op.TTL, op.Created)
		case models.Delete:
			g.cache.Delete(op.Key, op.Created)
		}
	}
	g.logger.Info("Apply operations from: %s. applyed %d operations", msg.Sender, len(msg.Operations))
}

func (g *Gossip) startMembershipGossip(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	g.logger.Info("Start sending membership gossips")
	for {
		select {
		case <-ticker.C:
			g.sendMembershipGossip()
		case <-ctx.Done():
			g.logger.Info("End sending membership gossips")
			return
		}
	}
}

func (g *Gossip) startDataGossip(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	g.logger.Info("Start sending data gossips")
	for {
		select {
		case <-ticker.C:
			g.sendDataGossip()
		case <-ctx.Done():
			g.logger.Info("End sending data gossips")
			return
		}
	}
}
func (g *Gossip) sendMembershipGossip() {
	member := g.membership.GetMember()
	g.logger.Info("Send member gossip to %s", member.Name)

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
	g.membership.UpdateSendedTime(member)
}

func (g *Gossip) sendDataGossip() {
	member := g.membership.GetMember()
	g.logger.Info("Send data gossip to %s", member.Name)
	allOp := g.cache.GetPendingOperations(member.LastIndex)
	if len(allOp) <= 0 {
		g.logger.Info("Nothing to send")
		return
	}
	operations := make([]models.Operation, len(allOp))
	maxIndex := 0
	for i, o := range allOp {
		operations = append(operations, models.Operation{o.Key, o.Value, o.Type, o.Ttl, o.Created})
		maxIndex = i

	}
	msg := models.NewOperationsGossip(g.nodeName, g.address, operations)

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
	g.membership.UpdateSendedTime(member)
	g.membership.UpdateMember(member, maxIndex)
}
