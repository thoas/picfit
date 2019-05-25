package store

import (
	"fmt"
	"time"

	"github.com/ulule/gokvstores"

	"github.com/thoas/picfit/logger"
)

type Store gokvstores.KVStore

const (
	dummyKVStoreType        = "dummy"
	redisKVStoreType        = "redis"
	redisClusterKVStoreType = "redis-cluster"
	cacheKVStoreType        = "cache"
)

// New returns a KVStore from config
func New(log logger.Logger, cfg *Config) (gokvstores.KVStore, error) {
	if cfg == nil {
		return gokvstores.DummyStore{}, nil
	}

	log.Debug("KVStore configured",
		logger.String("type", cfg.Type))

	switch cfg.Type {
	case dummyKVStoreType:
		return gokvstores.DummyStore{}, nil
	case redisClusterKVStoreType:
		redis := cfg.RedisCluster

		s, err := gokvstores.NewRedisClusterStore(&gokvstores.RedisClusterOptions{
			Addrs:    redis.Addrs,
			Password: redis.Password,
		}, time.Duration(redis.Expiration)*time.Second)
		if err != nil {
			return nil, err
		}

		return &kvstoreWrapper{s, cfg.Prefix}, nil
	case redisKVStoreType:
		redis := cfg.Redis

		s, err := gokvstores.NewRedisClientStore(&gokvstores.RedisClientOptions{
			Addr:     redis.Addr(),
			DB:       redis.DB,
			Password: redis.Password,
		}, time.Duration(redis.Expiration)*time.Second)
		if err != nil {
			return nil, err
		}

		return &kvstoreWrapper{s, cfg.Prefix}, nil
	case cacheKVStoreType:
		cache := cfg.Cache

		s, err := gokvstores.NewMemoryStore(
			time.Duration(cache.Expiration)*time.Second,
			time.Duration(cache.CleanupInterval)*time.Second)
		if err != nil {
			return nil, err
		}

		return &kvstoreWrapper{s, cfg.Prefix}, nil
	}

	return nil, fmt.Errorf("kvstore %s does not exist", cfg.Type)
}
