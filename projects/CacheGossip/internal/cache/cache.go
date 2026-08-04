package cache

import (
	"sync"
	"time"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

type cacheItem struct {
	value string
	ttl   time.Time
}

func NewCache() (*Cache, error) {
	return &Cache{sync.RWMutex{}, make(map[string]cacheItem)}, nil
}

func (c *Cache) Set(key, value string, ttl time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = cacheItem{value, ttl}
	return nil
}

func (c *Cache) Remove(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exist := c.items[key]

	if !exist {
		return "", false
	}
	if time.Now().After(item.ttl) {
		delete(c.items, key)
		return "", false
	}
	return item.value, true
}
