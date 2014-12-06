package kvstores

type KVStore interface {
	NewFromParams(params map[string]string) KVStore
	Connection() KVStoreConnection
	Close() error
}

type KVStoreConnection interface {
	Close() error
	Get(key string) string
	Set(key string, value string) error
}
