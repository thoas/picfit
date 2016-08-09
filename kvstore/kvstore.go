package kvstore

import (
	"fmt"

	"github.com/thoas/gokvstores"
	"github.com/thoas/picfit/config"
)

// NewKVStoreFromConfig returns a KVStore from config
func NewKVStoreFromConfig(cfg *config.Config) (gokvstores.KVStore, error) {
	if cfg.KVStore == nil {
		return &DummyKVStore{}, nil
	}

	section := cfg.KVStore

	switch section.Type {
	case "redis":
		host := section.Host

		password := section.Password

		db := section.Db

		port := section.Port

		return gokvstores.NewRedisKVStore(host, port, password, db), nil
	case "cache":
		if section.MaxEntries == 0 {
			section.MaxEntries = -1
		}

		return gokvstores.NewCacheKVStore(section.MaxEntries), nil

	}

	return nil, fmt.Errorf("kvstore %s does not exist", section.Type)
}
