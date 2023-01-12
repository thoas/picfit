package store

import (
	"context"
	"strings"
	"time"

	"github.com/ulule/gokvstores"
)

const (
	replicaErrorMessage = "READONLY You can't write against a read only replica."
)

type redisRoundRobinStore struct {
	kvstores []gokvstores.KVStore
}

func (k *redisRoundRobinStore) forEach(f func(kvstore gokvstores.KVStore) error) error {
	var err error
	for i := range k.kvstores {
		err = f(k.kvstores[i])
		if err == nil {
			return nil
		}
		if !strings.Contains(err.Error(), replicaErrorMessage) {
			return err
		}
	}

	return err
}

func (k *redisRoundRobinStore) Set(ctx context.Context, key string, value interface{}) error {
	return k.forEach(func(kvstore gokvstores.KVStore) error {
		return kvstore.Set(ctx, key, value)
	})
}

func (k *redisRoundRobinStore) Get(ctx context.Context, key string) (interface{}, error) {
	var (
		res interface{}
		err error
	)
	if err := k.forEach(func(kvstore gokvstores.KVStore) error {
		res, err = kvstore.Get(ctx, key)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return res, err
	}

	return res, nil
}

func (k *redisRoundRobinStore) Exists(ctx context.Context, keys ...string) (bool, error) {
	var (
		res bool
		err error
	)
	if err := k.forEach(func(kvstore gokvstores.KVStore) error {
		res, err = kvstore.Exists(ctx, keys...)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return res, err
	}

	return res, nil
}

func (k *redisRoundRobinStore) AppendSlice(ctx context.Context, key string, values ...interface{}) error {
	return k.forEach(func(kvstore gokvstores.KVStore) error {
		return kvstore.AppendSlice(ctx, key, values...)
	})
}

func (k *redisRoundRobinStore) GetSlice(ctx context.Context, key string) ([]interface{}, error) {
	var (
		res []interface{}
		err error
	)
	if err := k.forEach(func(kvstore gokvstores.KVStore) error {
		res, err = kvstore.GetSlice(ctx, key)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return res, err
	}

	return res, nil
}

func (k *redisRoundRobinStore) Delete(ctx context.Context, key string) error {
	return k.forEach(func(kvstore gokvstores.KVStore) error {
		return kvstore.Delete(ctx, key)
	})
}

func (k *redisRoundRobinStore) Close() error {
	return k.forEach(func(kvstore gokvstores.KVStore) error {
		return kvstore.Close()
	})
}

func (k *redisRoundRobinStore) DeleteMap(ctx context.Context, key string, fields ...string) error {
	return k.forEach(func(kvstore gokvstores.KVStore) error {
		return kvstore.DeleteMap(ctx, key, fields...)
	})
}

func (k *redisRoundRobinStore) Flush(ctx context.Context) error {
	return k.forEach(func(kvstore gokvstores.KVStore) error {
		return kvstore.Flush(ctx)
	})
}

func (k *redisRoundRobinStore) GetMaps(ctx context.Context, keys []string) (map[string]map[string]interface{}, error) {
	var (
		res map[string]map[string]interface{}
		err error
	)
	if err := k.forEach(func(kvstore gokvstores.KVStore) error {
		res, err = kvstore.GetMaps(ctx, keys)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return res, err
	}

	return res, nil
}

func (k *redisRoundRobinStore) SetMap(ctx context.Context, key string, value map[string]interface{}) error {
	return k.forEach(func(kvstore gokvstores.KVStore) error {
		return kvstore.SetMap(ctx, key, value)
	})
}

func (k *redisRoundRobinStore) GetMap(ctx context.Context, key string) (map[string]interface{}, error) {
	var (
		res map[string]interface{}
		err error
	)
	if err := k.forEach(func(kvstore gokvstores.KVStore) error {
		res, err = kvstore.GetMap(ctx, key)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return res, err
	}

	return res, nil
}

func (k *redisRoundRobinStore) Keys(ctx context.Context, pattern string) ([]interface{}, error) {
	var (
		res []interface{}
		err error
	)
	if err := k.forEach(func(kvstore gokvstores.KVStore) error {
		res, err = kvstore.Keys(ctx, pattern)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return res, err
	}

	return res, nil
}

func (k *redisRoundRobinStore) MGet(ctx context.Context, keys []string) (map[string]interface{}, error) {
	var (
		res map[string]interface{}
		err error
	)
	if err := k.forEach(func(kvstore gokvstores.KVStore) error {
		res, err = kvstore.MGet(ctx, keys)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return res, err
	}

	return res, nil
}

func (k *redisRoundRobinStore) SetMaps(ctx context.Context, maps map[string]map[string]interface{}) error {
	return k.forEach(func(kvstore gokvstores.KVStore) error {
		return kvstore.SetMaps(ctx, maps)
	})
}

func (k *redisRoundRobinStore) SetSlice(ctx context.Context, key string, value []interface{}) error {
	return k.forEach(func(kvstore gokvstores.KVStore) error {
		return kvstore.SetSlice(ctx, key, value)
	})
}

func (k *redisRoundRobinStore) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return k.forEach(func(kvstore gokvstores.KVStore) error {
		return kvstore.SetWithExpiration(ctx, key, value, expiration)
	})
}
