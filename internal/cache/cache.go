package cache

import (
	"sync"
)

// Cache struct represents cache. Based on sync.Map
type Cache struct {
	data sync.Map
}

// InitCache return new instance of cache
func InitCache() *Cache {
	return &Cache{sync.Map{}}
}

// Add data if not exists in cache. Returns true if exists, else returns false
func (c *Cache) Add(key string, value any) {
	c.data.Store(key, value)
	return
}

// Get data from cache. Returns true if data exists, else returns false
func (c *Cache) Get(key string) (value any, ok bool) {
	value, ok = c.data.Load(key)
	return
}
