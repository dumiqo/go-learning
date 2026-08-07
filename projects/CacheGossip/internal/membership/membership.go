package membership

import (
	"CacheGossip/pkg/models"
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Membership struct {
	members map[string]member
	mu      sync.RWMutex
}

type member struct {
	Name               string
	Url                string
	Status             models.Status
	LastSeen, LastSend time.Time
	LastIndex          int
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

	member.LastSeen = time.Now().UTC()
	member.Status = models.Dead

	m.members[name] = member

	return true
}
func (m *Membership) IsKnownMember(name string) bool {
	_, exist := m.members[name]
	return exist

}
func (m *Membership) UpdateStatus(name string, time time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	member, exists := m.members[name]
	if !exists {
		return
	}

	member.LastSeen = time
	member.Status = models.NewStatus(time)

	m.members[name] = member
}
func (m *Membership) AddMember(
	name string,
	url string,
	lastSeen time.Time) {
	m.members[name] = member{name, url, models.NewStatus(lastSeen), lastSeen, time.Time{}, 0}
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
		m.members[node.Name] = updateMember(member, node)
	}

	return nil
}

func newMember(node models.NodeInfo) member {
	return member{node.Name, node.Address, node.Status, node.LastSeen, time.Time{}, 0}
}
func updateMember(m member, node models.NodeInfo) member {
	nm := member{node.Name, node.Address, node.Status, node.LastSeen, time.Time{}, m.LastIndex}

	nm.LastSend = m.LastSend
	return nm
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

func (m *Membership) GetMember() member {
	return m.lastSendMember()
}

func (m *Membership) randomMember() member {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := m.randomMemberKey()
	return m.members[key]
}

func (m *Membership) lastSendMember() member {
	m.mu.RLock()
	defer m.mu.RUnlock()
	keys := make([]string, 0, len(m.members))

	for _, v := range m.members {
		keys = append(keys, v.Name)
	}

	sort.Slice(keys, func(i, j int) bool {
		return m.members[keys[i]].LastSend.Before(m.members[keys[j]].LastSend)
	})

	return m.members[keys[0]]
}

func (m *Membership) UpdateMember(member member, lastSendIndex int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, exist := m.members[member.Name]
	if !exist {
		return
	}
	item.LastSend = time.Now().UTC()
	m.members[member.Name] = item
}

func (m *Membership) UpdateSendedTime(member member) {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, exist := m.members[member.Name]
	if !exist {
		return
	}
	item.LastSend = time.Now().UTC()
	m.members[member.Name] = item
}

func (m *Membership) randomMemberKey() string {
	keys := make([]string, 0, len(m.members))
	for k, _ := range m.members {
		keys = append(keys, k)
	}
	return randomItem(keys)
}

func (m *Membership) randomAliveMemberKey() string {
	keys := make([]string, 0, len(m.members))
	alive := make([]string, 0, len(m.members))
	for k, member := range m.members {
		if member.Status != models.Dead {
			alive = append(alive, k)
		}
		keys = append(keys, k)
	}
	if len(alive) == 0 {
		return randomItem(keys)
	}
	return randomItem(alive)
}

func randomItem(slice []string) string {
	i := rand.Intn(len(slice))
	return slice[i]
}

func NewMembership(hosts map[string]string) *Membership {
	members := make(map[string]member)
	for k, v := range hosts {
		members[k] = member{k, v, models.Suspect, time.Time{}, time.Time{}, 0}
	}
	return &Membership{members, sync.RWMutex{}}
}
