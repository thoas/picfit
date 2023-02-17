package storage

import (
	"fmt"
	"strings"

	"github.com/mitchellh/goamz/aws"
	"github.com/thoas/picfit/http"
	"github.com/ulule/gostorages"

	"github.com/thoas/picfit/logger"
)

const (
	DOs3StorageType     = "dos3"
	fsStorageType       = "fs"
	gcsStorageType      = "gcs"
	httpDOs3StorageType = "http+dos3"
	httpFSStorageType   = "http+fs"
	httpS3StorageType   = "http+s3"
	httpStoragePrefix   = "http+"
	s3StorageType       = "s3"
)

// New return destination and source storages from config
func New(log logger.Logger, cfg *Config) (gostorages.Storage, gostorages.Storage, error) {
	if cfg == nil {
		storage := &DummyStorage{}

		return storage, storage, nil
	}

	var (
		sourceStorage      gostorages.Storage
		destinationStorage gostorages.Storage
		err                error
	)

	if cfg.Source != nil {
		sourceStorage, err = newStorage(cfg.Source)
		if err != nil {
			return nil, nil, err
		}

		log.Debug("Source storage configured",
			logger.String("type", cfg.Source.Type))
	}

	if cfg.Destination == nil {
		log.Debug("Destination storage not set, source storage will be used",
			logger.String("type", cfg.Source.Type))

		return sourceStorage, sourceStorage, nil
	}

	destinationStorage, err = newStorage(cfg.Destination)
	if err != nil {
		return nil, nil, err
	}

	log.Debug("Destination storage configured",
		logger.String("type", cfg.Destination.Type))

	return sourceStorage, destinationStorage, nil
}

func newStorage(cfg *StorageConfig) (gostorages.Storage, error) {
	if cfg == nil {
		return &DummyStorage{}, nil
	}

	if strings.HasPrefix(cfg.Type, httpStoragePrefix) && cfg.BaseURL == "" {
		return nil, fmt.Errorf("HTTP Wrapper cannot be used without setting *base_url* in config")
	}

	switch cfg.Type {
	case httpS3StorageType:
		cfg.Type = s3StorageType

		storage, err := newStorage(cfg)
		if err != nil {
			return nil, err
		}

		return NewHTTPStorage(storage, http.NewClient()), nil
	case s3StorageType:
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
	case httpDOs3StorageType:
		cfg.Type = DOs3StorageType

		storage, err := newStorage(cfg)
		if err != nil {
			return nil, err
		}

		return NewHTTPStorage(storage, http.NewClient()), nil
	case DOs3StorageType:
		acl, ok := gostorages.ACLs[cfg.ACL]
		if !ok {
			return nil, fmt.Errorf("The ACL %s does not exist", cfg.ACL)
		}

		region, ok := GetDOs3Region(cfg.Region)
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
	case gcsStorageType:
		return gostorages.NewGCSStorage(
			cfg.SecretAccessKey,
			cfg.BucketName,
			cfg.Location,
			cfg.BaseURL,
			cfg.CacheControl)
	case fsStorageType:
		return gostorages.NewFileSystemStorage(cfg.Location, cfg.BaseURL), nil
	case httpFSStorageType:
		cfg.Type = fsStorageType

		storage, err := newStorage(cfg)
		if err != nil {
			return nil, err
		}

		return NewHTTPStorage(storage, http.NewClient()), nil
	}

	return nil, fmt.Errorf("storage %s does not exist", cfg.Type)
}
