package membership

import (
	"CacheGossip/pkg/models"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Membership struct {
	members map[string]member
	mu      sync.RWMutex
}

type member struct {
	Name     string
	Url      string
	Status   models.Status
	LastSeen time.Time
}
type MembershipStatus struct {
	Members []MemberStatus
}
type MemberStatus struct {
	Name     string
	Url      string
	Status   models.Status
	LastSeen time.Time
}

func (m *Membership) Status() *MembershipStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mStatus := make([]MemberStatus, 0, len(m.members))

	for _, v := range m.members {
		mStatus = append(mStatus, MemberStatus{v.Name, v.Url, v.Status, v.LastSeen})
	}

	return &MembershipStatus{mStatus}
}

func (m *Membership) Start(ctx context.Context) error {

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.updateStatus()
		case <-ctx.Done():
			return nil
		}
	}
}

func (m *Membership) updateStatus() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range m.members {
		v.Status = models.NewStatus(v.LastSeen)
		m.members[k] = v
	}
}
func (m *Membership) KillNode(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	member, exists := m.members[name]
	if !exists {
		return false
	}

	member.LastSeen = time.Now()
	member.Status = models.Dead

	m.members[name] = member

	return true
}
func (m *Membership) UpdateStatus(name string, time time.Time) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	member, exists := m.members[name]
	if !exists {
		return false
	}

	member.LastSeen = time
	member.Status = models.NewStatus(time)

	m.members[name] = member

	return true
}
func (m *Membership) AddMember(
	name string,
	url string,
	lastSeen time.Time) {
	m.members[name] = member{name, url, models.NewStatus(lastSeen), lastSeen}
}
func (m *Membership) Update(node models.NodeInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	member, exists := m.members[node.Name]
	if !exists {
		m.members[node.Name] = newMember(node)
		return nil
	}

	if member.Name == node.Name {
		if node.LastSeen.Before(member.LastSeen) {
			return fmt.Errorf("Cannot update node")
		}
		m.members[node.Name] = newMember(node)
	}

	return nil
}

func newMember(node models.NodeInfo) member {
	return member{node.Name, node.Address, node.Status, node.LastSeen}
}

func (m *member) toNodeInfo() models.NodeInfo {
	return models.NodeInfo{Name: m.Name, Status: m.Status, Address: m.Url, LastSeen: m.LastSeen}
}

func (m *Membership) GetNodeInfo() []models.NodeInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	nodes := make([]models.NodeInfo, 0, len(m.members))

	for _, m := range m.members {
		nodes = append(nodes, m.toNodeInfo())
	}

	return nodes
}

func (m *Membership) RandomMember() member {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := m.randomMemberKey()
	return m.members[key]
}

func (m *Membership) randomMemberKey() string {
	keys := make([]string, len(m.members))
	for k, _ := range m.members {
		keys = append(keys, k)
	}
	i := rand.Intn(len(keys))
	return keys[i]
}

func NewMembership(hosts map[string]string) *Membership {
	members := make(map[string]member)
	for k, v := range hosts {
		members[k] = member{k, v, models.Suspect, time.Time{}}
	}
	return &Membership{members, sync.RWMutex{}}
}
