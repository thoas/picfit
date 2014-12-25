package application

import (
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/thoas/kvstores"
	"github.com/thoas/picfit/http"
	"github.com/thoas/storages"
	"strconv"
)

type KVStoreParameter func(params map[string]string) (kvstores.KVStore, error)
type StorageParameter func(params map[string]string) (storages.Storage, error)

var CacheKVStoreParameter KVStoreParameter = func(params map[string]string) (kvstores.KVStore, error) {
	value, ok := params["max_entries"]

	var maxEntries int

	if !ok {
		maxEntries = -1
	} else {
		maxEntries, _ = strconv.Atoi(value)
	}

	return kvstores.NewCacheKVStore(maxEntries), nil
}

var RedisKVStoreParameter KVStoreParameter = func(params map[string]string) (kvstores.KVStore, error) {
	host := params["host"]

	password := params["password"]

	port, _ := strconv.Atoi(params["port"])

	db, _ := strconv.Atoi(params["db"])

	return kvstores.NewRedisKVStore(host, port, password, db), nil
}

var FileSystemStorageParameter StorageParameter = func(params map[string]string) (storages.Storage, error) {
	return storages.NewFileSystemStorage(params["location"], params["base_url"]), nil
}

var HTTPFileSystemStorageParameter StorageParameter = func(params map[string]string) (storages.Storage, error) {
	storage, err := FileSystemStorageParameter(params)

	if err != nil {
		return nil, err
	}

	return &http.HTTPStorage{storage}, nil
}

var HTTPS3StorageParameter StorageParameter = func(params map[string]string) (storages.Storage, error) {
	storage, err := S3StorageParameter(params)

	if err != nil {
		return nil, err
	}

	return &http.HTTPStorage{storage}, nil
}

var S3StorageParameter StorageParameter = func(params map[string]string) (storages.Storage, error) {

	ACL, ok := storages.ACLs[params["acl"]]

	if !ok {
		return nil, fmt.Errorf("The ACL %s does not exist", params["acl"])
	}

	Region, ok := aws.Regions[params["region"]]

	if !ok {
		return nil, fmt.Errorf("The Region %s does not exist", params["region"])
	}

	return storages.NewS3Storage(
		params["access_key_id"],
		params["secret_access_key"],
		params["bucket_name"],
		params["location"],
		Region,
		ACL,
		params["base_url"],
	), nil
}
