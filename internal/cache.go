// Package pokecache is for making a cache for the results
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
	sync.Mutex
	data map[string]cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		data: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, data []byte) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       data,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Lock()
	defer c.Unlock()

	entry, ok := c.data[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	tick := time.Tick(interval)
	for range tick {
		c.Lock()
		for key, v := range c.data {
			if time.Since(v.createdAt) > interval {
				delete(c.data, key)
			}
		}
		c.Unlock()
	}
}
