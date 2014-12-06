package kvstores

type KVStore interface {
	NewFromParams(params map[string]string) KVStore
	Connection() KVStoreConnection
	Close() error
}

type KVStoreConnection interface {
	Close() error
	Get(key string) string
	Delete(key string) error
	Exists(key string) bool
	Set(key string, value string) error
}
