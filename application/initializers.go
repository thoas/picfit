package application

import (
	"fmt"
	"github.com/jmoiron/jsonq"
	"github.com/thoas/picfit/dummy"
	"github.com/thoas/picfit/image"
	"github.com/thoas/storages"
)

type Initializer func(jq *jsonq.JsonQuery) error

var Initializers = []Initializer{
	KVStoreInitializer,
	StorageInitializer,
	ShardInitializer,
	BasicInitializer,
}

var BasicInitializer Initializer = func(jq *jsonq.JsonQuery) error {
	format, _ := jq.String("format")

	if format != "" {
		App.Format = format
		App.ContentType = image.ContentTypes[format]
	} else {
		App.Format = DefaultFormat
		App.ContentType = DefaultContentType
	}

	App.SecretKey, _ = jq.String("secret_key")

	return nil
}

var ShardInitializer Initializer = func(jq *jsonq.JsonQuery) error {
	width, err := jq.Int("shard", "width")

	if err != nil {
		width = DefaultShardWidth
	}

	depth, err := jq.Int("shard", "depth")

	if err != nil {
		depth = DefaultShardDepth
	}

	App.Shard = Shard{Width: width, Depth: depth}

	return nil
}

var KVStoreInitializer Initializer = func(jq *jsonq.JsonQuery) error {
	_, err := jq.Object("kvstore")

	if err != nil {
		App.KVStore = &dummy.DummyKVStore{}

		return nil
	}

	key, err := jq.String("kvstore", "type")

	if err != nil {
		return err
	}

	parameter, ok := KVStores[key]

	if !ok {
		return fmt.Errorf("KVStore %s does not exist", key)
	}

	config, err := jq.Object("kvstore")

	if err != nil {
		return err
	}

	store, err := parameter(mapInterfaceToMapString(config))

	if err != nil {
		return err
	}

	App.KVStore = store

	return nil
}

func getStorageFromConfig(key string, jq *jsonq.JsonQuery) (storages.Storage, error) {
	storageType, err := jq.String("storage", key, "type")

	parameter, ok := Storages[storageType]

	if !ok {
		return nil, fmt.Errorf("Storage %s does not exist", key)
	}

	config, err := jq.Object("storage", key)

	if err != nil {
		return nil, err
	}

	storage, err := parameter(mapInterfaceToMapString(config))

	if err != nil {
		return nil, err
	}

	return storage, err
}

var StorageInitializer Initializer = func(jq *jsonq.JsonQuery) error {
	_, err := jq.Object("storage")

	if err != nil {
		App.SourceStorage = &dummy.DummyStorage{}
		App.DestStorage = &dummy.DummyStorage{}

		return nil
	}

	sourceStorage, err := getStorageFromConfig("src", jq)

	if err != nil {
		return err
	}

	App.SourceStorage = sourceStorage

	destStorage, err := getStorageFromConfig("dest", jq)

	if err != nil {
		App.DestStorage = sourceStorage
	}

	App.DestStorage = destStorage

	return nil
}
