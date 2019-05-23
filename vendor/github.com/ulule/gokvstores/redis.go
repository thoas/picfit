package gokvstores

import (
	"net"
	"time"

	conv "github.com/cstockton/go-conv"
	redis "gopkg.in/redis.v5"
)

// ----------------------------------------------------------------------------
// Client
// ----------------------------------------------------------------------------

// RedisClient is an interface thats allows to use Redis cluster or a redis single client seamlessly.
type RedisClient interface {
	Ping() *redis.StatusCmd
	Exists(key string) *redis.BoolCmd
	Del(keys ...string) *redis.IntCmd
	FlushDb() *redis.StatusCmd
	Close() error
	Process(cmd redis.Cmder) error
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	MGet(keys ...string) *redis.SliceCmd
	HDel(key string, fields ...string) *redis.IntCmd
	HGetAll(key string) *redis.StringStringMapCmd
	HMSet(key string, fields map[string]string) *redis.StatusCmd
	SMembers(key string) *redis.StringSliceCmd
	SAdd(key string, members ...interface{}) *redis.IntCmd
	Keys(pattern string) *redis.StringSliceCmd
	Pipeline() *redis.Pipeline
}

// RedisPipeline is a struct which contains an opend redis pipeline transaction
type RedisPipeline struct {
	pipeline *redis.Pipeline
}

// RedisClientOptions are Redis client options.
type RedisClientOptions struct {
	Network            string
	Addr               string
	Dialer             func() (net.Conn, error)
	Password           string
	DB                 int
	MaxRetries         int
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
	ReadOnly           bool
}

// RedisClusterOptions are Redis cluster options.
type RedisClusterOptions struct {
	Addrs              []string
	MaxRedirects       int
	ReadOnly           bool
	RouteByLatency     bool
	Password           string
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
	IdleTimeout        time.Duration
	IdleCheckFrequency time.Duration
}

// ----------------------------------------------------------------------------
// Store
// ----------------------------------------------------------------------------

// RedisStore is the Redis implementation of KVStore.
type RedisStore struct {
	client     RedisClient
	expiration time.Duration
}

// Get returns value for the given key.
func (r *RedisStore) Get(key string) (interface{}, error) {
	cmd := redis.NewCmd("get", key)

	if err := r.client.Process(cmd); err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	return cmd.Val(), cmd.Err()
}

// MGet returns map of key, value for a list of keys.
func (r *RedisStore) MGet(keys []string) (map[string]interface{}, error) {
	values, err := r.client.MGet(keys...).Result()

	newValues := make(map[string]interface{}, len(keys))

	for k, v := range keys {
		value := values[k]
		if err != nil {
			return nil, err
		}

		newValues[v] = value
	}
	return newValues, nil
}

// Set sets the value for the given key.
func (r *RedisStore) Set(key string, value interface{}) error {
	return r.client.Set(key, value, r.expiration).Err()
}

// SetWithExpiration sets the value for the given key.
func (r *RedisStore) SetWithExpiration(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(key, value, expiration).Err()
}

// GetMap returns map for the given key.
func (r *RedisStore) GetMap(key string) (map[string]interface{}, error) {
	values, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, nil
	}

	newValues := make(map[string]interface{}, len(values))
	for k, v := range values {
		newValues[k] = v
	}

	return newValues, nil
}

// SetMap sets map for the given key.
func (r *RedisStore) SetMap(key string, values map[string]interface{}) error {
	newValues := make(map[string]string, len(values))

	for k, v := range values {
		val, err := conv.String(v)
		if err != nil {
			return err
		}

		newValues[k] = val
	}

	return r.client.HMSet(key, newValues).Err()
}

// DeleteMap removes the specified fields from the map stored at key.
func (r *RedisStore) DeleteMap(key string, fields ...string) error {
	return r.client.HDel(key, fields...).Err()
}

// GetSlice returns slice for the given key.
func (r *RedisStore) GetSlice(key string) ([]interface{}, error) {
	values, err := r.client.SMembers(key).Result()
	if err != nil {
		return nil, err
	}

	if len(values) == 0 {
		return nil, nil
	}

	newValues := make([]interface{}, len(values))
	for i := range values {
		newValues[i] = values[i]
	}

	return newValues, nil
}

// SetSlice sets map for the given key.
func (r *RedisStore) SetSlice(key string, values []interface{}) error {
	for _, v := range values {
		if v != nil {
			if err := r.client.SAdd(key, v).Err(); err != nil {
				return err
			}
		}
	}

	return nil
}

// AppendSlice appends values to the given slice.
func (r *RedisStore) AppendSlice(key string, values ...interface{}) error {
	return r.SetSlice(key, values)
}

// Exists checks key existence.
func (r *RedisStore) Exists(key string) (bool, error) {
	cmd := r.client.Exists(key)
	return cmd.Val(), cmd.Err()
}

// Delete deletes key.
func (r *RedisStore) Delete(key string) error {
	return r.client.Del(key).Err()
}

// Keys returns all keys matching pattern.
func (r *RedisStore) Keys(pattern string) ([]interface{}, error) {
	values, err := r.client.Keys(pattern).Result()

	if len(values) == 0 {
		return nil, err
	}

	newValues := make([]interface{}, len(values))

	for k, v := range values {
		newValues[k] = v
	}

	return newValues, err
}

// Flush flushes the current database.
func (r *RedisStore) Flush() error {
	return r.client.FlushDb().Err()
}

// Close closes the client connection.
func (r *RedisStore) Close() error {
	return r.client.Close()
}

// NewRedisClientStore returns Redis client instance of KVStore.
func NewRedisClientStore(options *RedisClientOptions, expiration time.Duration) (KVStore, error) {
	opts := &redis.Options{
		Network:            options.Network,
		Addr:               options.Addr,
		Dialer:             options.Dialer,
		Password:           options.Password,
		DB:                 options.DB,
		MaxRetries:         options.MaxRetries,
		DialTimeout:        options.DialTimeout,
		ReadTimeout:        options.ReadTimeout,
		WriteTimeout:       options.WriteTimeout,
		PoolSize:           options.PoolSize,
		PoolTimeout:        options.PoolTimeout,
		IdleTimeout:        options.IdleTimeout,
		IdleCheckFrequency: options.IdleCheckFrequency,
		ReadOnly:           options.ReadOnly,
	}

	client := redis.NewClient(opts)

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &RedisStore{
		client:     client,
		expiration: expiration,
	}, nil
}

// NewRedisClusterStore returns Redis cluster client instance of KVStore.
func NewRedisClusterStore(options *RedisClusterOptions, expiration time.Duration) (KVStore, error) {
	opts := &redis.ClusterOptions{
		Addrs:              options.Addrs,
		MaxRedirects:       options.MaxRedirects,
		ReadOnly:           options.ReadOnly,
		RouteByLatency:     options.RouteByLatency,
		Password:           options.Password,
		DialTimeout:        options.DialTimeout,
		ReadTimeout:        options.ReadTimeout,
		WriteTimeout:       options.WriteTimeout,
		PoolSize:           options.PoolSize,
		PoolTimeout:        options.PoolTimeout,
		IdleTimeout:        options.IdleTimeout,
		IdleCheckFrequency: options.IdleCheckFrequency,
	}

	client := redis.NewClusterClient(opts)

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &RedisStore{
		client:     client,
		expiration: expiration,
	}, nil
}

// Pipeline uses pipeline as a Redis client to execute multiple calls at once
func (r *RedisStore) Pipeline(f func(r *RedisStore) error) ([]redis.Cmder, error) {
	pipe := r.client.Pipeline()

	redisPipeline := RedisPipeline{
		pipeline: pipe,
	}

	store := &RedisStore{
		client:     redisPipeline,
		expiration: r.expiration,
	}

	err := f(store)
	if err != nil {
		return nil, err
	}

	cmds, err := pipe.Exec()
	return cmds, err
}

// GetMaps returns maps for the given keys.
func (r *RedisStore) GetMaps(keys []string) (map[string]map[string]interface{}, error) {
	commands, err := r.Pipeline(func(r *RedisStore) error {
		for _, key := range keys {
			r.client.HGetAll(key)
		}
		return nil

	})
	if err != nil {
		return nil, err
	}

	newValues := make(map[string]map[string]interface{}, len(keys))

	for i, key := range keys {
		cmd := commands[i]
		values, _ := cmd.(*redis.StringStringMapCmd).Result()
		if values != nil {
			valueMap := make(map[string]interface{}, len(values))
			for k, v := range values {
				valueMap[k] = v
			}

			newValues[key] = valueMap
		} else {
			newValues[key] = nil
		}
	}

	return newValues, nil
}

// SetMaps sets the given maps.
func (r *RedisStore) SetMaps(maps map[string]map[string]interface{}) error {
	_, err := r.Pipeline(func(r *RedisStore) error {
		for k, v := range maps {
			r.SetMap(k, v)
		}
		return nil

	})
	return err
}

// Pipeline returns Redis pipeline
func (r RedisPipeline) Pipeline() *redis.Pipeline {
	return r.pipeline
}

// Ping implements RedisClient Ping for pipeline
func (r RedisPipeline) Ping() *redis.StatusCmd {
	return r.pipeline.Ping()
}

// Exists implements RedisClient Exists for pipeline
func (r RedisPipeline) Exists(key string) *redis.BoolCmd {
	return r.pipeline.Exists(key)
}

// Del implements RedisClient Del for pipeline
func (r RedisPipeline) Del(keys ...string) *redis.IntCmd {
	return r.pipeline.Del(keys...)
}

// FlushDb implements RedisClient FlushDb for pipeline
func (r RedisPipeline) FlushDb() *redis.StatusCmd {
	return r.pipeline.FlushDb()
}

// Close implements RedisClient Close for pipeline
func (r RedisPipeline) Close() error {
	return r.pipeline.Close()
}

// Process implements RedisClient Process for pipeline
func (r RedisPipeline) Process(cmd redis.Cmder) error {
	return r.pipeline.Process(cmd)
}

// Get implements RedisClient Get for pipeline
func (r RedisPipeline) Get(key string) *redis.StringCmd {
	return r.pipeline.Get(key)
}

// MGet implements RedisClient MGet for pipeline
func (r RedisPipeline) MGet(keys ...string) *redis.SliceCmd {
	return r.pipeline.MGet(keys...)
}

// Set implements RedisClient Set for pipeline
func (r RedisPipeline) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.pipeline.Set(key, value, expiration)
}

// HDel implements RedisClient HDel for pipeline
func (r RedisPipeline) HDel(key string, fields ...string) *redis.IntCmd {
	return r.pipeline.HDel(key, fields...)
}

// HGetAll implements RedisClient HGetAll for pipeline
func (r RedisPipeline) HGetAll(key string) *redis.StringStringMapCmd {
	return r.pipeline.HGetAll(key)
}

// HMSet implements RedisClient HMSet for pipeline
func (r RedisPipeline) HMSet(key string, fields map[string]string) *redis.StatusCmd {
	return r.pipeline.HMSet(key, fields)
}

// SMembers implements RedisClient SMembers for pipeline
func (r RedisPipeline) SMembers(key string) *redis.StringSliceCmd {
	return r.pipeline.SMembers(key)
}

// SAdd implements RedisClient SAdd for pipeline
func (r RedisPipeline) SAdd(key string, members ...interface{}) *redis.IntCmd {
	return r.pipeline.SAdd(key, members...)
}

// Keys implements RedisClient Keys for pipeline
func (r RedisPipeline) Keys(pattern string) *redis.StringSliceCmd {
	return r.pipeline.Keys(pattern)
}
