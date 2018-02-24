package kvstore

import "fmt"

// Config is a struct to represent a key/value store (redis, cache)
type Config struct {
	Type       string
	Prefix     string
	MaxEntries int
	Redis      RedisKVStore
	Cache      CacheKVStore
}

type RedisKVStore struct {
	Host       string
	Port       int
	Password   string
	DB         int
	Expiration int
}

func (r RedisKVStore) Addr() string {
	return fmt.Sprint(r.Host, ":", r.Port)
}

type CacheKVStore struct {
	Expiration      int
	CleanupInterval int `mapstructure:"cleanup_interval"`
}
