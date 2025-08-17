package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	data     map[string]cacheEntry
	interval time.Duration
	mu       *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	var newCache Cache
	newCache.data = make(map[string]cacheEntry)
	newCache.interval = interval
	newCache.mu = &sync.Mutex{}
	newCache.reapLoop()
	return &newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ce, ok := c.data[key]
	if !ok {
		return nil, false
	}
	return ce.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)

	go func() {
		for t := range ticker.C {
			c.mu.Lock()
			for k, ce := range c.data {
				if t.Sub(ce.createdAt) > c.interval {
					delete(c.data, k)
				}
			}
			c.mu.Unlock()
		}
	}()
}
