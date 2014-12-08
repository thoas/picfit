package dummy

import (
	"github.com/thoas/kvstores"
)

type DummyKVStore struct {
}

func (k *DummyKVStore) Connection() kvstores.KVStoreConnection {
	return &DummyKVStoreConnection{}
}

func (k *DummyKVStore) Close() error {
	return nil
}

type DummyKVStoreConnection struct {
}

func (k *DummyKVStoreConnection) Close() error {
	return nil
}

func (k *DummyKVStoreConnection) Get(key string) string {
	return ""
}

func (k *DummyKVStoreConnection) Delete(key string) error {
	return nil
}

func (k *DummyKVStoreConnection) Exists(key string) bool {
	return false
}

func (k *DummyKVStoreConnection) Set(key string, value string) error {
	return nil
}
