package store

import (
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

func (k *kvstoreWrapper) Set(key string, value interface{}) error {
	return k.KVStore.Set(k.prefixed(key), value)
}

func (k *kvstoreWrapper) Get(key string) (interface{}, error) {
	return k.KVStore.Get(k.prefixed(key))
}

func (k *kvstoreWrapper) Exists(key string) (bool, error) {
	return k.KVStore.Exists(k.prefixed(key))
}

func (k *kvstoreWrapper) AppendSlice(key string, values ...interface{}) error {
	return k.KVStore.AppendSlice(k.prefixed(key), values...)
}

func (k *kvstoreWrapper) GetSlice(key string) ([]interface{}, error) {
	return k.KVStore.GetSlice(k.prefixed(key))
}

func (k *kvstoreWrapper) Delete(key string) error {
	return k.KVStore.Delete(k.prefixed(key))
}
