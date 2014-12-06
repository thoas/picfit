package kvstores

type KVStore interface {
	Connect(params map[string]string)
	Connection() KVStoreConnection
	Close() error
}

type KVStoreConnection interface {
	Close() error
	Get(key string) string
	Set(key string, value string) error
}
