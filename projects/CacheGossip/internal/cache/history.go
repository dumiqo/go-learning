package cache

import (
	"CacheGossip/pkg/models"
	"sync"
	"time"
)

type History struct {
	mu         sync.RWMutex
	operations []Operation
}

func NewHistory() *History {
	return &History{sync.RWMutex{}, make([]Operation, 0)}
}

type Operation struct {
	Type         models.OperationType
	Key, Value   string
	Ttl, Created time.Time
}

func (h *History) Set(key, value string, ttl, time time.Time) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.operations = append(h.operations, Operation{models.Set, key, value, ttl, time})
	return nil
}

func (h *History) Delete(key string, time time.Time) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.operations = append(h.operations, Operation{models.Delete, key, "", time, time})
	return nil
}

func (h *History) Get(from int) []Operation {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.operations) <= from {
		return make([]Operation, 0)
	}

	ops := h.operations[from:]
	m := make(map[string]Operation)
	for _, v := range ops {
		op, exist := m[v.Key]
		if !exist || op.Created.Before(v.Created) {
			m[v.Key] = v
		}
	}
	result := make([]Operation, 0, len(m))

	for _, v := range m {
		result = append(result, v)
	}

	return result
}
