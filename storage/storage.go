package storage

import (
	"fmt"
	"strings"

	"github.com/mitchellh/goamz/aws"
	"github.com/thoas/gostorages"
	"github.com/thoas/picfit/config"
)

// NewStoragesFromConfig return destination and source storages from config
func NewStoragesFromConfig(cfg *config.Config) (gostorages.Storage, gostorages.Storage, error) {
	if cfg.Storage == nil {
		storage := &DummyStorage{}

		return storage, storage, nil
	}

	var sourceStorage gostorages.Storage
	var destinationStorage gostorages.Storage
	var err error

	if cfg.Storage.Src != nil {
		sourceStorage, err = NewStorageFromConfig(cfg.Storage.Src)

		if err != nil {
			return nil, nil, err
		}
	}

	if cfg.Storage.Dst == nil {
		return sourceStorage, sourceStorage, nil
	}

	destinationStorage, err = NewStorageFromConfig(cfg.Storage.Dst)

	if err != nil {
		return nil, nil, err
	}

	return sourceStorage, destinationStorage, nil
}

// NewStorageFromConfig returns a Storage from config
func NewStorageFromConfig(cfg *config.Storage) (gostorages.Storage, error) {
	if cfg == nil {
		return &DummyStorage{}, nil
	}

	if strings.HasPrefix(cfg.Type, "http+") && cfg.BaseURL == "" {
		return nil, fmt.Errorf("HTTP Wrapper cannot be used without setting *base_url* in config")
	}

	switch cfg.Type {
	case "http+s3":
		cfg.Type = "s3"

		storage, err := NewStorageFromConfig(cfg)

		if err != nil {
			return nil, err
		}

		return &HTTPStorage{storage, ""}, nil
	case "s3":
		acl, ok := gostorages.ACLs[cfg.ACL]

		if !ok {
			return nil, fmt.Errorf("The ACL %s does not exist", cfg.ACL)
		}

		region, ok := aws.Regions[cfg.Region]

		if !ok {
			return nil, fmt.Errorf("The Region %s does not exist", cfg.Region)
		}

		return gostorages.NewS3Storage(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			cfg.BucketName,
			cfg.Location,
			region,
			acl,
			cfg.BaseURL,
		), nil
	case "fs":
		return gostorages.NewFileSystemStorage(cfg.Location, cfg.BaseURL), nil
	case "http+fs":
		cfg.Type = "fs"

		storage, err := NewStorageFromConfig(cfg)

		if err != nil {
			return nil, err
		}

		return &HTTPStorage{storage, ""}, nil
	}

	return nil, fmt.Errorf("storage %s does not exist", cfg.Type)
}
