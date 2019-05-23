package gokvstores

import (
	"sort"
	"time"

	conv "github.com/cstockton/go-conv"
)

// KVStore is the KV store interface.
type KVStore interface {
	// Get returns value for the given key.
	Get(key string) (interface{}, error)

	// MGet returns map of key, value for a list of keys.
	MGet(keys []string) (map[string]interface{}, error)

	// Set sets value for the given key.
	Set(key string, value interface{}) error

	// SetWithExpiration sets the value for the given key for a specified duration.
	SetWithExpiration(key string, value interface{}, expiration time.Duration) error

	// GetMap returns map for the given key.
	GetMap(key string) (map[string]interface{}, error)

	// GetMaps returns maps for the given keys.
	GetMaps(keys []string) (map[string]map[string]interface{}, error)

	// SetMap sets map for the given key.
	SetMap(key string, value map[string]interface{}) error

	// SetMaps sets the given maps.
	SetMaps(maps map[string]map[string]interface{}) error

	// DeleteMap removes the specified fields from the map stored at key.
	DeleteMap(key string, fields ...string) error

	// GetSlice returns slice for the given key.
	GetSlice(key string) ([]interface{}, error)

	// SetSlice sets slice for the given key.
	SetSlice(key string, value []interface{}) error

	// AppendSlice appends values to an existing slice.
	// If key does not exist, creates slice.
	AppendSlice(key string, values ...interface{}) error

	// Exists checks if the given key exists.
	Exists(key string) (bool, error)

	// Delete deletes the given key.
	Delete(key string) error

	// Flush flushes the store.
	Flush() error

	// Return all keys matching pattern
	Keys(pattern string) ([]interface{}, error)

	// Close closes the connection to the store.
	Close() error
}

func stringSlice(values []interface{}) ([]string, error) {
	converted := []string{}

	for _, v := range values {
		if v != nil {
			val, err := conv.String(v)
			if err != nil {
				return nil, err
			}

			converted = append(converted, val)
		}
	}

	sort.Strings(converted)

	return converted, nil
}
