package kvstore

import (
	"github.com/thoas/gokvstores"
)

type DummyKVStore struct {
}

func (k *DummyKVStore) Connection() gokvstores.KVStoreConnection {
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

func (k *DummyKVStoreConnection) Flush() error {
	return nil
}

func (k *DummyKVStoreConnection) Get(key string) interface{} {
	return ""
}

func (k *DummyKVStoreConnection) Delete(key string) error {
	return nil
}

func (k *DummyKVStoreConnection) Exists(key string) bool {
	return false
}

func (k *DummyKVStoreConnection) Set(key string, value interface{}) error {
	return nil
}

func (k *DummyKVStoreConnection) Append(key string, value interface{}) error {
	return nil
}

func (k *DummyKVStoreConnection) SetAdd(key string, value interface{}) error {
	return nil
}

func (k *DummyKVStoreConnection) SetMembers(key string) []interface{} {
	return nil
}
