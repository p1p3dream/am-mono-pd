package val

import (
	"sync"
)

// Cache can be used to store any type of key/value pairs in a goroutine-safe way.
type Cache[K comparable, V any] struct {
	mu    sync.Mutex
	items map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]V),
	}
}

func (c *Cache[K, V]) Set(k K, v V) {
	c.mu.Lock()
	c.items[k] = v
	c.mu.Unlock()
}

func (c *Cache[K, V]) Get(k K) V {
	c.mu.Lock()
	v := c.items[k]
	c.mu.Unlock()
	return v
}

func (c *Cache[K, V]) Select(k K) (V, bool) {
	c.mu.Lock()
	v, ok := c.items[k]
	c.mu.Unlock()
	return v, ok
}

func (c *Cache[K, V]) Has(k K) bool {
	c.mu.Lock()
	_, ok := c.items[k]
	c.mu.Unlock()
	return ok
}

func (c *Cache[K, V]) SetGet(k K, v V) V {
	c.mu.Lock()
	c.items[k] = v
	c.mu.Unlock()
	return v
}

func (c *Cache[K, V]) SetGetOnce(k K, v V) V {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[k]; ok {
		return item
	}

	c.items[k] = v
	return v
}

func (c *Cache[K, V]) Pop(k K) (V, bool) {
	c.mu.Lock()
	v, ok := c.items[k]
	if ok {
		delete(c.items, k)
	}
	c.mu.Unlock()
	return v, ok
}

func (c *Cache[K, V]) Del(k K) {
	c.mu.Lock()
	delete(c.items, k)
	c.mu.Unlock()
}
