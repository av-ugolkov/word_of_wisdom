package cache

import (
	"sync"

	"github.com/google/uuid"
)

type Cache struct {
	mx    sync.Mutex
	cache map[uuid.UUID]struct{}
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[uuid.UUID]struct{}),
	}
}

func (c *Cache) Add(key uuid.UUID) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.cache[key] = struct{}{}
}

func (c *Cache) Get(key uuid.UUID) bool {
	c.mx.Lock()
	defer c.mx.Unlock()
	_, ok := c.cache[key]
	return ok
}

// Delete - delete key from cache
func (c *Cache) Delete(key uuid.UUID) {
	c.mx.Lock()
	defer c.mx.Unlock()
	delete(c.cache, key)
}
