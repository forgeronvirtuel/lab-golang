package main

import (
	"fmt"
	"sync"
)

type LockedCache struct {
	*sync.Mutex
	data map[string]string
}

// NewLockedCache creates and returns a *LockedCache with prefill values.
func NewLockedCache(data map[string]string) *LockedCache {
	if data == nil {
		data = map[string]string{}
	}
	return &LockedCache{
		Mutex: &sync.Mutex{},
		data:  data,
	}
}

// Add add a key / value pair into the cache. If a value
// already exist, replace it. It is threadsafe.
func (c *LockedCache) Add(key, value string) {
	c.Lock()
	c.data[key] = value
	c.Unlock()
}

// GetWithStatus return the value into the cache and a boolean that
// indicates if a value was found. It is threadsafe. If `c`
// is nil, act as an empty cache.
func (c *LockedCache) GetWithStatus(key string) (string, bool) {
	c.Lock()
	if c == nil {
		return "", false
	}
	v, ok := c.data[key]
	c.Unlock()
	return v, ok
}

// Get return the value into the cache. It is threadsafe.
// If `c` is nil, act as an empty cache.
func (c *LockedCache) Get(key string) string {
	if c == nil {
		return ""
	}
	c.Lock()
	v, _ := c.data[key]
	c.Unlock()
	return v
}

func main() {
	// Nil cache usable as an empty cache (at least to get data)
	var cache *LockedCache
	fmt.Println("Getting some key value:", cache.Get("some key"))

	// Creating a cache and add a value
	cache = NewLockedCache(nil)
	cache.Add("some key", "some value")
	fmt.Println("Getting some key value:", cache.Get("some key"))

	// We have access to the internal structure
	fmt.Println(cache.Mutex)

	// Oops, deadlock !
	cache.Lock()
	fmt.Println("Getting some key value:", cache.Get("some key"))
	cache.Unlock()
}
