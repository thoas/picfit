package gokvstores

import "time"

// DummyStore is a noop store (caching disabled).
type DummyStore struct{}

// Get returns value for the given key.
func (DummyStore) Get(key string) (interface{}, error) {
	return nil, nil
}

// MGet returns map of key, value for a list of keys.
func (DummyStore) MGet(keys []string) (map[string]interface{}, error) {
	return nil, nil
}

// Set sets value for the given key.
func (DummyStore) Set(key string, value interface{}) error {
	return nil
}

// SetWithExpiration sets the value for the given key for a specified duration.
func (DummyStore) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	return nil
}

// GetMap returns map for the given key.
func (DummyStore) GetMap(key string) (map[string]interface{}, error) {
	return nil, nil
}

// GetMaps returns maps for the given keys.
func (DummyStore) GetMaps(keys []string) (map[string]map[string]interface{}, error) {
	return nil, nil
}

// SetMap sets map for the given key.
func (DummyStore) SetMap(key string, value map[string]interface{}) error {
	return nil
}

// SetMaps sets the given maps.
func (DummyStore) SetMaps(maps map[string]map[string]interface{}) error { return nil }

// DeleteMap removes the specified fields from the map stored at key.
func (DummyStore) DeleteMap(key string, fields ...string) error { return nil }

// GetSlice returns slice for the given key.
func (DummyStore) GetSlice(key string) ([]interface{}, error) {
	return nil, nil
}

// SetSlice sets slice for the given key.
func (DummyStore) SetSlice(key string, value []interface{}) error {
	return nil
}

// AppendSlice appends values to an existing slice.
// If key does not exist, creates slice.
func (DummyStore) AppendSlice(key string, values ...interface{}) error {
	return nil
}

// Exists checks if the given key exists.
func (DummyStore) Exists(key string) (bool, error) {
	return false, nil
}

// Delete deletes the given key.
func (DummyStore) Delete(key string) error {
	return nil
}

// Keys returns all keys matching pattern
func (DummyStore) Keys(pattern string) ([]interface{}, error) {
	return nil, nil
}

// Flush flushes the store.
func (DummyStore) Flush() error {
	return nil
}

// Close closes the connection to the store.
func (DummyStore) Close() error {
	return nil
}
