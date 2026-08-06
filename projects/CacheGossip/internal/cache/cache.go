package cache

import (
	"CacheGossip/pkg/logger"
	"CacheGossip/pkg/models"
	"context"
	"sync"
	"time"
)

type Cache struct {
	mu    sync.RWMutex
	log   *logger.Logger
	items map[string]cacheItem

	miss, hit int
}

type CacheStat struct {
	Count, Miss, Hit int
}

type cacheItem struct {
	value string
	ttl   time.Time
}

func (c *Cache) autoClean(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.cleanExpired()
		case <-ctx.Done():
			return // ← корректное завершение
		}
	}
}
func (c *Cache) cleanExpired() {
	c.log.Info("Start autoclean")
	c.mu.RLock()
	now := time.Now()
	toDelete := make([]string, 0, 100)

	for k, v := range c.items {
		if now.After(v.ttl) {
			toDelete = append(toDelete, k)
		}
		if len(toDelete) >= 100 {
			break
		}
	}
	c.mu.RUnlock()

	c.log.Info("Items to delete: %d", len(toDelete))
	if len(toDelete) > 0 {
		c.mu.Lock()
		for _, k := range toDelete {
			delete(c.items, k)
		}
		c.mu.Unlock()
	}
	c.log.Info("End autoclean")
}

func NewCache(logger *logger.Logger) (*Cache, error) {
	return &Cache{sync.RWMutex{}, logger, make(map[string]cacheItem), 0, 0}, nil
}

func (c *Cache) Start(cnx context.Context) {
	c.autoClean(cnx)
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
		c.miss++
		return "", false
	}
	if time.Now().After(item.ttl) {
		delete(c.items, key)
		return "", false
	}
	c.hit++
	return item.value, true
}

func (c *Cache) Stat() CacheStat {
	return CacheStat{len(c.items), c.miss, c.hit}
}
func (c *Cache) GetPendingOperations() []models.PendingOperation {
	return CacheStat{len(c.items), c.miss, c.hit}
}
