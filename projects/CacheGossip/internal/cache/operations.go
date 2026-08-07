package cache

import (
	"CacheGossip/pkg/models"
	"sync"
	"time"
)

type OperationBuffer struct {
	mu         sync.RWMutex
	operations map[string]Operation
}

func NewOperations() *OperationBuffer {
	return &OperationBuffer{sync.RWMutex{}, make(map[string]Operation)}
}

type Operation struct {
	Type         models.OperationType
	Key, Value   string
	Ttl, Created time.Time
}

func (p *OperationBuffer) Set(key, value string, ttl, time time.Time) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	v, exist := p.operations[key]
	if !exist || v.Created.Before(time) {
		p.operations[key] = Operation{models.Set, key, value, ttl, time}
	}
	return nil
}

func (p *OperationBuffer) Delete(key string, time time.Time) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	v, exist := p.operations[key]
	if exist && v.Created.Before(time) {
		p.operations[key] = Operation{models.Delete, key, "", time, time}
	}
	return nil
}

func (p *OperationBuffer) Get() []Operation {
	p.mu.Lock()
	defer p.mu.Unlock()
	operations := make([]Operation, 0, len(p.operations))
	for _, o := range p.operations {
		operations = append(operations, o)
	}
	clear(p.operations)
	return operations
}
