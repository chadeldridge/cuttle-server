package db

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	CACHE_DEFAULT_EXPIRE = 5 * time.Minute
	CACHE_DEFAULT_PURGE  = 10 * time.Minute
)

var ErrKeyNotFound = fmt.Errorf("key not found")

// Cache is a wrapper around the go-cache library.
type Cache struct {
	instance *cache.Cache
}

// DefaultCache returns a new Cache instance with default expiration and purge times.
func DefaultCache() *Cache {
	return &Cache{instance: cache.New(CACHE_DEFAULT_EXPIRE, CACHE_DEFAULT_PURGE)}
}

// NewCache returns a new Cache instance with the specified expiration and purge times.
func NewCache(expire, purge time.Duration) *Cache {
	return &Cache{instance: cache.New(expire, purge)}
}

// Set the key to the value with the instance's DefaultExpiration.
func (c *Cache) Set(key string, value interface{}) {
	c.instance.Set(key, value, cache.DefaultExpiration)
}

// Set the key to the value with the specified expiration time.
func (c *Cache) SetWithExpire(key string, value interface{}, expire time.Duration) {
	c.instance.Set(key, value, expire)
}

// Get the value for the key. Returns value and true if found.
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.instance.Get(key)
}

// Get the value for the key and refresh the expiration time to the instance's DefaultExpiration.
func (c *Cache) GetAndRefresh(key string) (interface{}, bool) {
	v, f := c.instance.Get(key)

	// If the key is not found, return the value and false.
	if !f {
		return v, f
	}

	// Set the key with the value and the instance's DefaultExpiration.
	c.instance.Set(key, v, cache.DefaultExpiration)
	return v, f
}

// Delete the key from the cache.
func (c *Cache) Delete(key string) {
	c.instance.Delete(key)
}

// Refresh the expiration time for the key to the instance's DefaultExpiration.
func (c *Cache) Refresh(key string) error {
	v, f := c.instance.Get(key)
	if !f {
		return fmt.Errorf("Cache.Refresh: %s", ErrKeyNotFound)
	}

	// Set the key with the value and the instance's DefaultExpiration.
	c.instance.Set(key, v, cache.DefaultExpiration)
	return nil
}
