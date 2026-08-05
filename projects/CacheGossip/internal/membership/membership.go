package membership

import (
	"CacheGossip/pkg/models"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Membership struct {
	members map[string]Member
}

type Member struct {
	Name    string
	Url     string
	Status  models.Status
	History []history
}
type history struct {
	uuid uuid.UUID
	time time.Time
}

func (m *Membership) RandomMember() Member {
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
	members := make(map[string]Member)
	for k, v := range hosts {
		members[k] = Member{k, v, models.Suspect, make([]history, 10)}
	}
	return &Membership{members}
}
