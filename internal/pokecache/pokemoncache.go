package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Pokemon struct {
	Name string `json:"name"`
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}

	go c.reapLoop()

	return c
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.entries {
			if (entry.createdAt.Add(c.interval)).Before(time.Now()) {
				delete(c.entries, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Add(cacheKey string, val []byte) {
	c.mu.Lock()
	c.entries[cacheKey] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	c.mu.Unlock()

}

func (c *Cache) Get(cacheKey string) ([]byte, bool) {
	c.mu.Lock()
	for key, entry := range c.entries {
		if key == cacheKey {
			c.mu.Unlock()
			return entry.val, true
		}
	}
	c.mu.Unlock()
	return nil, false
}
