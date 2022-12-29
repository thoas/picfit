package store

import (
	"context"
	"fmt"

	"github.com/ulule/gokvstores"
)

type kvstoreWrapper struct {
	gokvstores.KVStore

	Prefix string
}

func (k *kvstoreWrapper) prefixed(key string) string {
	return fmt.Sprint(k.Prefix, key)
}

func (k *kvstoreWrapper) Set(ctx context.Context, key string, value interface{}) error {
	return k.KVStore.Set(ctx, k.prefixed(key), value)
}

func (k *kvstoreWrapper) Get(ctx context.Context, key string) (interface{}, error) {
	return k.KVStore.Get(ctx, k.prefixed(key))
}

func (k *kvstoreWrapper) Exists(ctx context.Context, keys ...string) (bool, error) {
	newkeys := make([]string, len(keys))
	for i := range keys {
		newkeys[i] = k.prefixed(keys[i])
	}
	return k.KVStore.Exists(ctx, newkeys...)
}

func (k *kvstoreWrapper) AppendSlice(ctx context.Context, key string, values ...interface{}) error {
	return k.KVStore.AppendSlice(ctx, k.prefixed(key), values...)
}

func (k *kvstoreWrapper) GetSlice(ctx context.Context, key string) ([]interface{}, error) {
	return k.KVStore.GetSlice(ctx, k.prefixed(key))
}

func (k *kvstoreWrapper) Delete(ctx context.Context, key string) error {
	return k.KVStore.Delete(ctx, k.prefixed(key))
}
