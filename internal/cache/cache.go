package cache

import (
	"fmt"
	"log"
	"sync"
)

type InMemory struct {
	m map[string][]byte
	//  A read/write mutex allows all readers to access the map at the same time, but a writer will lock out everyone else
	lock sync.RWMutex
}

// New initializes a new map to use for storing cache
func New() *InMemory {
	m := make(map[string][]byte)
	c := InMemory{m: m}
	return &c
}

// Set Adds a new order to the cache
func (c *InMemory) Set(key string, value []byte) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.m[key]; !ok {
		c.m[key] = value
		log.Print("Added to cache id ", key, string(value))

		return nil
	}
	return fmt.Errorf("Already exists")
}

// Get reads data from cache and returns found value for the provided key and error
func (c *InMemory) Get(key string) (value []byte, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if mapValue, ok := c.m[key]; ok {
		value = append(value[:0], mapValue...)
		return value, nil
	}

	return value, fmt.Errorf("No such value")
}

// IsEmpty checks if cache stores any data
func (c *InMemory) IsEmpty() (ans bool) {
	if len(c.m) == 0 {
		return true
	}
	return false
}
