package store

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/ulule/gokvstores"
	"go.uber.org/zap"

	"github.com/thoas/picfit/logger"
)

type Store gokvstores.KVStore

const (
	cacheKVStoreType        = "cache"
	dummyKVStoreType        = "dummy"
	redisClusterKVStoreType = "redis-cluster"
	redisKVStoreType        = "redis"
	redisRoundRobinType     = "redis-roundrobin"
)

func parseRedisURL(redisURL string) (*gokvstores.RedisClientOptions, error) {
	opts := &gokvstores.RedisClientOptions{}
	u, err := url.Parse(redisURL)
	if err != nil {
		return nil, err
	}
	h, p := getHostPortWithDefaults(u)
	opts.Addr = net.JoinHostPort(h, p)
	q := u.Query()
	opts.Password = q.Get("password")
	rawdb := q.Get("db")
	if rawdb != "" {
		db, err := strconv.Atoi(rawdb)
		if err != nil {
			return nil, err
		}
		opts.DB = db
	} else {
		opts.DB = 0
	}
	return opts, nil
}

func getHostPortWithDefaults(u *url.URL) (string, string) {
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}
	return host, port
}

// New returns a KVStore from config
func New(ctx context.Context, log *zap.Logger, cfg *Config) (gokvstores.KVStore, error) {
	if cfg == nil {
		return gokvstores.DummyStore{}, nil
	}

	log.Debug("KVStore configured",
		logger.String("type", cfg.Type))

	switch cfg.Type {
	case dummyKVStoreType:
		return gokvstores.DummyStore{}, nil
	case redisRoundRobinType:
		redis := cfg.RedisRoundRobin

		kvstores := make([]gokvstores.KVStore, len(redis.Addrs))
		for i := range redis.Addrs {
			redisOptions, err := parseRedisURL(redis.Addrs[i])
			if err != nil {
				return nil, err
			}
			s, err := gokvstores.NewRedisClientStore(ctx, redisOptions, time.Duration(redis.Expiration)*time.Second)
			if err != nil {
				return nil, err
			}

			kvstores[i] = s
		}

		return &kvstoreWrapper{&redisRoundRobinStore{kvstores}, cfg.Prefix}, nil
	case redisClusterKVStoreType:
		redis := cfg.RedisCluster

		s, err := gokvstores.NewRedisClusterStore(ctx, &gokvstores.RedisClusterOptions{
			Addrs:    redis.Addrs,
			Password: redis.Password,
		}, time.Duration(redis.Expiration)*time.Second)
		if err != nil {
			return nil, err
		}

		return &kvstoreWrapper{s, cfg.Prefix}, nil
	case redisKVStoreType:
		redis := cfg.Redis

		s, err := gokvstores.NewRedisClientStore(ctx, &gokvstores.RedisClientOptions{
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
