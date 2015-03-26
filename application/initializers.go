package application

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/jmoiron/jsonq"
	"github.com/thoas/gostorages"
	"github.com/thoas/picfit/dummy"
	"github.com/thoas/picfit/engines"
	"github.com/thoas/picfit/util"
)

type Initializer func(jq *jsonq.JsonQuery, app *Application) error

var Initializers = []Initializer{
	KVStoreInitializer,
	StorageInitializer,
	ShardInitializer,
	BasicInitializer,
	SentryInitializer,
}

var KVStores = map[string]KVStoreParameter{
	"redis": RedisKVStoreParameter,
	"cache": CacheKVStoreParameter,
}

var Storages = map[string]StorageParameter{
	"http+s3": HTTPS3StorageParameter,
	"s3":      S3StorageParameter,
	"http+fs": HTTPFileSystemStorageParameter,
	"fs":      FileSystemStorageParameter,
}

var SentryInitializer Initializer = func(jq *jsonq.JsonQuery, app *Application) error {
	dsn, err := jq.String("sentry", "dsn")

	if err != nil {
		return nil
	}

	results, err := jq.Object("sentry", "tags")

	var tags map[string]string

	if err != nil {
		tags = map[string]string{}
	} else {
		tags = util.MapInterfaceToMapString(results)
	}

	client, err := raven.NewClient(dsn, tags)

	if err != nil {
		return err
	}

	app.Raven = client

	return nil
}

var BasicInitializer Initializer = func(jq *jsonq.JsonQuery, app *Application) error {
	f, _ := jq.String("format")

	var format string

	if f != "" {
		format = f
	} else {
		format = DefaultFormat
	}

	app.SecretKey, _ = jq.String("secret_key")
	app.Engine = engines.NewGoImageEngine(format)

	return nil
}

var ShardInitializer Initializer = func(jq *jsonq.JsonQuery, app *Application) error {
	width, err := jq.Int("shard", "width")

	if err != nil {
		width = DefaultShardWidth
	}

	depth, err := jq.Int("shard", "depth")

	if err != nil {
		depth = DefaultShardDepth
	}

	app.Shard = Shard{Width: width, Depth: depth}

	return nil
}

var KVStoreInitializer Initializer = func(jq *jsonq.JsonQuery, app *Application) error {
	_, err := jq.Object("kvstore")

	if err != nil {
		app.KVStore = &dummy.DummyKVStore{}

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

	params := util.MapInterfaceToMapString(config)
	store, err := parameter(params)

	if err != nil {
		return err
	}

	app.Prefix = params["prefix"]
	app.KVStore = store

	return nil
}

func getStorageFromConfig(key string, jq *jsonq.JsonQuery) (gostorages.Storage, error) {
	storageType, err := jq.String("storage", key, "type")

	parameter, ok := Storages[storageType]

	if !ok {
		return nil, fmt.Errorf("Storage %s does not exist", key)
	}

	config, err := jq.Object("storage", key)

	if err != nil {
		return nil, err
	}

	storage, err := parameter(util.MapInterfaceToMapString(config))

	if err != nil {
		return nil, err
	}

	return storage, err
}

var StorageInitializer Initializer = func(jq *jsonq.JsonQuery, app *Application) error {
	_, err := jq.Object("storage")

	if err != nil {
		app.SourceStorage = &dummy.DummyStorage{}
		app.DestStorage = &dummy.DummyStorage{}

		return nil
	}

	sourceStorage, err := getStorageFromConfig("src", jq)

	if err != nil {
		return err
	}

	app.SourceStorage = sourceStorage

	destStorage, err := getStorageFromConfig("dst", jq)

	if err != nil {
		app.DestStorage = sourceStorage
	} else {
		app.DestStorage = destStorage
	}

	return nil
}
