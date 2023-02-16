package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/goamz/aws"
	"github.com/thoas/picfit/http"
	"github.com/ulule/gostorages"
	fsstorage "github.com/ulule/gostorages/fs"
	gcstorage "github.com/ulule/gostorages/gcs"
	s3storage "github.com/ulule/gostorages/s3"
	"go.uber.org/zap"

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

type Storage struct {
	gostorages.Storage
	StorageConfig
}

// New return destination and source storages from config
func New(log *zap.Logger, cfg *Config) (*Storage, *Storage, error) {
	if cfg == nil {
		storage := &Storage{Storage: &DummyStorage{}}

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

		return &Storage{
				Storage:       sourceStorage,
				StorageConfig: *cfg.Source,
			}, &Storage{
				Storage:       sourceStorage,
				StorageConfig: *cfg.Source,
			}, nil
	}

	destinationStorage, err = newStorage(cfg.Destination)
	if err != nil {
		return nil, nil, err
	}

	log.Debug("Destination storage configured",
		logger.String("type", cfg.Destination.Type))

	return &Storage{
			Storage:       sourceStorage,
			StorageConfig: *cfg.Source,
		}, &Storage{
			Storage:       destinationStorage,
			StorageConfig: *cfg.Destination,
		}, nil
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
		region, ok := aws.Regions[cfg.Region]
		if !ok {
			return nil, fmt.Errorf("the region %s does not exist", cfg.Region)
		}

		return s3storage.NewStorage(s3storage.Config{
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
			Region:          region.Name,
			Bucket:          cfg.BucketName,
		})
	case httpDOs3StorageType:
		cfg.Type = DOs3StorageType

		storage, err := newStorage(cfg)
		if err != nil {
			return nil, err
		}

		return NewHTTPStorage(storage, http.NewClient()), nil
	case DOs3StorageType:
		region, ok := GetDOs3Region(cfg.Region)
		if !ok {
			return nil, fmt.Errorf("the region %s does not exist", cfg.Region)
		}
		return s3storage.NewStorage(s3storage.Config{
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
			Region:          region.Name,
			Bucket:          cfg.BucketName,
		})
	case gcsStorageType:
		return gcstorage.NewStorage(context.Background(), cfg.SecretAccessKey, cfg.BucketName)
	case fsStorageType:
		return fsstorage.NewStorage(fsstorage.Config{Root: cfg.Location}), nil
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
