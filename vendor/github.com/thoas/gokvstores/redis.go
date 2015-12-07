package gokvstores

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

func newPool(host string, port int, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func NewRedisKVStore(host string, port int, password string, db int) KVStore {
	return &RedisKVStore{Pool: newPool(host, port, password, db)}
}

type RedisKVStore struct {
	Pool *redis.Pool
}

func (k *RedisKVStore) Close() error {
	return k.Pool.Close()
}

func (k *RedisKVStore) Connection() KVStoreConnection {
	return &RedisKVStoreConnection{Connection: k.Pool.Get()}
}

type RedisKVStoreConnection struct {
	Connection redis.Conn
}

func (k *RedisKVStoreConnection) Close() error {
	return k.Connection.Close()
}

func (k *RedisKVStoreConnection) Flush() error {
	return k.Connection.Flush()
}

func (k *RedisKVStoreConnection) Get(key string) interface{} {
	reply, err := k.Connection.Do("GET", key)

	if err != nil {
		return nil
	}

	return reply
}

func (k *RedisKVStoreConnection) Exists(key string) bool {
	exists, err := redis.Bool(k.Connection.Do("EXISTS", key))

	if err != nil {
		return false
	}

	return exists
}

func (k *RedisKVStoreConnection) Set(key string, value interface{}) error {
	_, err := k.Connection.Do("SET", key, value)

	if err != nil {
		return err
	}

	return nil
}

func (k *RedisKVStoreConnection) Append(key string, value interface{}) error {
	_, err := k.Connection.Do("APPEND", key, value)

	if err != nil {
		return err
	}

	return nil
}

func (k *RedisKVStoreConnection) SetAdd(key string, value interface{}) error {
	_, err := k.Connection.Do("SADD", key, value)

	if err != nil {
		return err
	}

	return nil
}

func (k *RedisKVStoreConnection) SetMembers(key string) []interface{} {
	reply, err := k.Connection.Do("SMEMBERS", key)

	if err != nil {
		return nil
	}

	values, err := Values(reply)
	if err != nil {
		return nil
	}
	if len(values) == 0 {
		return nil
	}

	return values
}

func (k *RedisKVStoreConnection) Delete(key string) error {
	_, err := k.Connection.Do("DEL", key)

	if err != nil {
		return err
	}

	return nil
}
