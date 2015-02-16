package application

import (
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/thoas/gokvstores"
	"github.com/thoas/gostorages"
	"github.com/thoas/picfit/http"
	"strconv"
)

type KVStoreParameter func(params map[string]string) (gokvstores.KVStore, error)
type StorageParameter func(params map[string]string) (gostorages.Storage, error)

var CacheKVStoreParameter KVStoreParameter = func(params map[string]string) (gokvstores.KVStore, error) {
	value, ok := params["max_entries"]

	var maxEntries int

	if !ok {
		maxEntries = -1
	} else {
		maxEntries, _ = strconv.Atoi(value)
	}

	return gokvstores.NewCacheKVStore(maxEntries), nil
}

var RedisKVStoreParameter KVStoreParameter = func(params map[string]string) (gokvstores.KVStore, error) {
	host := params["host"]

	password := params["password"]

	port, _ := strconv.Atoi(params["port"])

	db, _ := strconv.Atoi(params["db"])

	return gokvstores.NewRedisKVStore(host, port, password, db), nil
}

var FileSystemStorageParameter StorageParameter = func(params map[string]string) (gostorages.Storage, error) {
	return gostorages.NewFileSystemStorage(params["location"], params["base_url"]), nil
}

var HTTPFileSystemStorageParameter StorageParameter = func(params map[string]string) (gostorages.Storage, error) {
	storage, err := FileSystemStorageParameter(params)

	if err != nil {
		return nil, err
	}

	if _, ok := params["base_url"]; !ok {
		return nil, fmt.Errorf("You can't use the http wrapper without setting *base_url* in your config file")
	}

	return &http.HTTPStorage{storage}, nil
}

var HTTPS3StorageParameter StorageParameter = func(params map[string]string) (gostorages.Storage, error) {
	storage, err := S3StorageParameter(params)

	if err != nil {
		return nil, err
	}

	if _, ok := params["base_url"]; !ok {
		return nil, fmt.Errorf("You can't use the http wrapper without setting *base_url* in your config file")
	}

	return &http.HTTPStorage{storage}, nil
}

var S3StorageParameter StorageParameter = func(params map[string]string) (gostorages.Storage, error) {

	ACL, ok := gostorages.ACLs[params["acl"]]

	if !ok {
		return nil, fmt.Errorf("The ACL %s does not exist", params["acl"])
	}

	Region, ok := aws.Regions[params["region"]]

	if !ok {
		return nil, fmt.Errorf("The Region %s does not exist", params["region"])
	}

	return gostorages.NewS3Storage(
		params["access_key_id"],
		params["secret_access_key"],
		params["bucket_name"],
		params["location"],
		Region,
		ACL,
		params["base_url"],
	), nil
}
