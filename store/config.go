package store

import "fmt"

// Config is a struct to represent a key/value store (redis, cache)
type Config struct {
	Cache           CacheConfig `mapstructure:"cache"`
	Prefix          string
	Redis           RedisConfig           `mapstructure:"redis"`
	RedisCluster    RedisClusterConfig    `mapstructure:"redis-cluster"`
	RedisRoundRobin RedisRoundRobinConfig `mapstructure:"redis-roundrobin"`
	Type            string
}

type RedisConfig struct {
	DB         int
	Expiration int
	Host       string
	Password   string
	Port       int
}

type RedisClusterConfig struct {
	Addrs      []string
	Expiration int
	Password   string
}

type RedisRoundRobinConfig struct {
	Addrs      []string
	Expiration int
	Password   string
}

func (r RedisConfig) Addr() string {
	return fmt.Sprint(r.Host, ":", r.Port)
}

type CacheConfig struct {
	CleanupInterval int `mapstructure:"cleanup_interval"`
	Expiration      int
}
