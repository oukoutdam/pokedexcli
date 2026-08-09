package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu            sync.RWMutex
	cacheEntryMap map[string]cacheEntry
	interval      time.Duration
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cacheEntryMap: make(map[string]cacheEntry),
		interval:      interval,
	}

	go c.reapLoop()

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheEntryMap[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cacheEntryMap[key]
	if !ok {
		return []byte{}, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for key, value := range c.cacheEntryMap {
		if now.Sub(value.createdAt) >= c.interval {
			delete(c.cacheEntryMap, key)
		}
	}
}
