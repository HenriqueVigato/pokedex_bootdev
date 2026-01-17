// Package pokecache is for making a cache for the results
package pokecache

import (
	"fmt"
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
	fmt.Println("Um cache foi criado")
	return cache
}

func (c *Cache) Add(key string, data []byte) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       data,
	}
	fmt.Println("Um cache foi adicionado")
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Lock()
	defer c.Unlock()
	fmt.Println("Um cache foi buscado")

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
				fmt.Println("Um cache foi excluido")
				delete(c.data, key)
			}
		}
		c.Unlock()
	}
}
