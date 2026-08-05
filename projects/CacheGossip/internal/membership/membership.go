package membership

import (
	"CacheGossip/pkg/models"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Membership struct {
	members map[string]member
}

type member struct {
	name    string
	url     string
	status  models.Status
	history []history
}
type history struct {
	uuid uuid.UUID
	time time.Time
}

func (m *Membership) autoNotifyMember() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		key := m.randomMemberKey()
		member := m.members[key]
	}
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
		members[k] = member{k, v, models.Suspect, make([]history, 10)}
	}
	m := Membership{members}
	go m.autoNotifyMember()

	return &m
}
