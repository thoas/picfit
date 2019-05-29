package store

import "fmt"

// Config is a struct to represent a key/value store (redis, cache)
type Config struct {
	Type         string
	Prefix       string
	Postgres     PostgresConfig     `mapstructure:"postgres"`
	Redis        RedisConfig        `mapstructure:"redis"`
	Cache        CacheConfig        `mapstructure:"cache"`
	RedisCluster RedisClusterConfig `mapstructure:"redis-cluster"`
}

type PostgresConfig struct {
	WriteDb string `mapstructure:"write_db"`
	ReadDb  string `mapstructure:"read_db"`
}

type RedisConfig struct {
	Host       string
	Port       int
	Password   string
	DB         int
	Expiration int
}

type RedisClusterConfig struct {
	Expiration int
	Password   string
	Addrs      []string
}

func (r RedisConfig) Addr() string {
	return fmt.Sprint(r.Host, ":", r.Port)
}

type CacheConfig struct {
	Expiration      int
	CleanupInterval int `mapstructure:"cleanup_interval"`
}
