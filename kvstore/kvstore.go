package kvstore

import (
	"fmt"
	"time"

	"github.com/ulule/gokvstores"
)

const (
	dummyKVStoreType = "dummy"
	redisKVStoreType = "redis"
	cacheKVStoreType = "cache"
)

// New returns a KVStore from config
func New(cfg *Config) (gokvstores.KVStore, error) {
	if cfg == nil {
		return gokvstores.DummyStore{}, nil
	}

	switch cfg.Type {
	case dummyKVStoreType:
		return gokvstores.DummyStore{}, nil
	case redisKVStoreType:
		redis := cfg.Redis

		return gokvstores.NewRedisClientStore(&gokvstores.RedisClientOptions{
			Addr:     redis.Addr(),
			DB:       redis.DB,
			Password: redis.Password,
		}, time.Duration(redis.Expiration)*time.Second)
	case cacheKVStoreType:
		cache := cfg.Cache

		return gokvstores.NewMemoryStore(
			time.Duration(cache.Expiration)*time.Second,
			time.Duration(cache.CleanupInterval)*time.Second)
	}

	return nil, fmt.Errorf("kvstore %s does not exist", cfg.Type)
}
