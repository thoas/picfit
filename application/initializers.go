package application

import (
	"errors"
	"fmt"
	"github.com/jmoiron/jsonq"
	"github.com/thoas/picfit/image"
	"strings"
)

type Initializer func(key string, jq *jsonq.JsonQuery) error

var Initializers = map[string]Initializer{
	"kvstore":  KVStoreInitializer,
	"storage":  StorageInitializer,
	"shard":    ShardInitializer,
	"format":   FormatInitializer,
	"base_url": BaseURLInitializer,
}

var BaseURLInitializer Initializer = func(key string, jq *jsonq.JsonQuery) error {
	App.BaseURL = strings.TrimSuffix(key, "/")

	return nil
}

var FormatInitializer Initializer = func(format string, jq *jsonq.JsonQuery) error {
	if format != "" {
		App.Format = format
		App.ContentType = image.ContentTypes[format]
	} else {
		App.Format = DefaultFormat
		App.ContentType = DefaultContentType
	}

	return nil
}

var ShardInitializer Initializer = func(key string, jq *jsonq.JsonQuery) error {
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

var KVStoreInitializer Initializer = func(key string, jq *jsonq.JsonQuery) error {
	store, ok := KVStores[key]

	if !ok {
		return errors.New(fmt.Sprintf("KVStore %s does not exist", key))
	}

	config, err := jq.Object(key)

	if err != nil {
		return err
	}

	store = store.NewFromParams(mapInterfaceToMapString(config))

	App.KVStore = store

	return nil
}

var StorageInitializer Initializer = func(key string, jq *jsonq.JsonQuery) error {
	storage, ok := Storages[key]

	if !ok {
		return errors.New(fmt.Sprintf("Storage %s does not exist", key))
	}

	config, err := jq.Object(key)

	if err != nil {
		return err
	}

	storage, err = storage.NewFromParams(mapInterfaceToMapString(config))

	App.Storage = storage

	return nil
}
