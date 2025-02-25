package utils

import (
	"sync"
	"time"
)

type cacheItem struct {
	value     string
	timestamp time.Time
}

type Cache struct {
	items map[string]cacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]cacheItem),
		ttl:   time.Hour,
	}

	go cache.cleanup()
	return cache
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:     value,
		timestamp: time.Now(),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return "", false
	}

	if time.Since(item.timestamp) > c.ttl {
		return "", false
	}

	return item.value, true
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		for key, item := range c.items {
			if time.Since(item.timestamp) > c.ttl {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
