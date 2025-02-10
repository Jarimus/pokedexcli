package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
	mu       *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	e := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Lock()
	c.entries[key] = e
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {

	// Start a ticker to check for expired entries
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		// on tick, check for expired entries and delete them
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > c.interval {
				c.mu.Lock()
				delete(c.entries, key)
				c.mu.Unlock()
			}
		}
	}

}
