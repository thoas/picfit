package kvstores

type DummyKVStore struct {
}

func (k *DummyKVStore) Connect(params map[string]string) {
}

func (k *DummyKVStore) Close() error {
	return nil
}

type DummyKVStoreConnection struct {
}

func (k *DummyKVStore) Connection() KVStoreConnection {
	return &DummyKVStoreConnection{}
}

func (k *DummyKVStoreConnection) Close() error {
	return nil
}

func (k *DummyKVStoreConnection) Get(key string) string {
	return ""
}

func (k *DummyKVStoreConnection) Set(key string, value string) error {
	return nil
}
