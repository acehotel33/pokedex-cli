package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mux     sync.RWMutex
}

func NewCache(interval time.Duration) *Cache {

	newCache := Cache{
		entries: map[string]cacheEntry{},
	}

	go newCache.reapLoop(interval)
	return &newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	if val, ok := c.entries[key]; ok {
		return val.val, true
	} else {
		return []byte{}, false
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mux.Lock()
		for url, entry := range c.entries {
			if time.Since(entry.createdAt) > interval {
				delete(c.entries, url)
			}
		}
		c.mux.Unlock()
	}

}
