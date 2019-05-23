package gokvstores

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// MemoryStore is the in-memory implementation of KVStore.
type MemoryStore struct {
	cache           *cache.Cache
	expiration      time.Duration
	cleanupInterval time.Duration
}

// Get returns item from the cache.
func (c *MemoryStore) Get(key string) (interface{}, error) {
	item, _ := c.cache.Get(key)
	return item, nil
}

// MGet returns map of key, value for a list of keys.
func (c *MemoryStore) MGet(keys []string) (map[string]interface{}, error) {
	results := make(map[string]interface{}, len(keys))
	for _, key := range keys {
		item, _ := c.Get(key)
		results[key] = item
	}
	return results, nil
}

// Set sets value in the cache.
func (c *MemoryStore) Set(key string, value interface{}) error {
	c.cache.Set(key, value, c.expiration)
	return nil
}

// SetWithExpiration sets the value for the given key for a specified duration.
func (c *MemoryStore) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	c.cache.Set(key, value, expiration)
	return nil
}

// GetMap returns map for the given key.
func (c *MemoryStore) GetMap(key string) (map[string]interface{}, error) {
	if v, found := c.cache.Get(key); found {
		return v.(map[string]interface{}), nil
	}
	return nil, nil
}

// GetMaps returns maps for the given keys.
func (c *MemoryStore) GetMaps(keys []string) (map[string]map[string]interface{}, error) {
	values := make(map[string]map[string]interface{}, len(keys))
	for _, v := range keys {
		value, _ := c.GetMap(v)
		if value != nil {
			values[v] = value
		}
	}

	return values, nil
}

// SetMap sets a map for the given key.
func (c *MemoryStore) SetMap(key string, value map[string]interface{}) error {
	c.cache.Set(key, value, c.expiration)
	return nil
}

// SetMaps sets the given maps.
func (c *MemoryStore) SetMaps(maps map[string]map[string]interface{}) error {
	for k, v := range maps {
		c.SetMap(k, v)
	}
	return nil
}

// DeleteMap removes the specified fields from the map stored at key.
func (c *MemoryStore) DeleteMap(key string, fields ...string) error {
	m, err := c.GetMap(key)
	if err != nil {
		return err
	}

	for _, field := range fields {
		delete(m, field)
	}

	return c.SetMap(key, m)
}

// GetSlice returns slice for the given key.
func (c *MemoryStore) GetSlice(key string) ([]interface{}, error) {
	if v, found := c.cache.Get(key); found {
		return v.([]interface{}), nil
	}
	return nil, nil
}

// SetSlice sets slice for the given key.
func (c *MemoryStore) SetSlice(key string, value []interface{}) error {
	c.cache.Set(key, value, c.expiration)
	return nil
}

// AppendSlice appends values to the given slice.
func (c *MemoryStore) AppendSlice(key string, values ...interface{}) error {
	items, err := c.GetSlice(key)
	if err != nil {
		return err
	}

	if items == nil {
		return c.SetSlice(key, values)
	}

	for _, item := range values {
		items = append(items, item)
	}

	return c.cache.Replace(key, items, c.expiration)
}

// Close does nothing for this backend.
func (c *MemoryStore) Close() error {
	return nil
}

// Flush removes all items from the cache.
func (c *MemoryStore) Flush() error {
	c.cache.Flush()
	return nil
}

// Delete deletes the given key.
func (c *MemoryStore) Delete(key string) error {
	c.cache.Delete(key)
	return nil
}

// Keys returns all keys matching pattern
func (c *MemoryStore) Keys(pattern string) ([]interface{}, error) {
	return nil, nil
}

// Exists checks if the given key exists.
func (c *MemoryStore) Exists(key string) (bool, error) {
	if _, exists := c.cache.Get(key); exists {
		return true, nil
	}
	return false, nil
}

// NewMemoryStore returns in-memory KVStore.
func NewMemoryStore(expiration time.Duration, cleanupInterval time.Duration) (KVStore, error) {
	return &MemoryStore{
		cache:           cache.New(expiration, cleanupInterval),
		expiration:      time.Duration(expiration) * time.Second,
		cleanupInterval: cleanupInterval,
	}, nil
}
