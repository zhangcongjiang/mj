package cache

import (
	"fmt"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	cache map[string]interface{}
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]interface{}),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, found := c.cache[key]
	return val, found
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = value
}

func main() {
	cache := NewCache()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	if val, found := cache.Get("key1"); found {
		fmt.Println("Value for key1:", val)
	} else {
		fmt.Println("Key1 not found in cache")
	}

	if val, found := cache.Get("key3"); found {
		fmt.Println("Value for key3:", val)
	} else {
		fmt.Println("Key3 not found in cache")
	}
}
