package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path"
	"strings"

	"github.com/thoas/picfit/http"
	"github.com/ulule/gostorages"
	fsstorage "github.com/ulule/gostorages/fs"
	gcstorage "github.com/ulule/gostorages/gcs"
	s3storage "github.com/ulule/gostorages/s3"
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

// Storage wraps gostorages.Storage.
type Storage struct {
	storage gostorages.Storage
	cfg     StorageConfig
}

// URL returns the filepath prefixed with BaseURL from storage.
func (s *Storage) URL(filepath string) string {
	if s.cfg.BaseURL != "" {
		if _, ok := s.storage.(*fsstorage.Storage); ok || s.cfg.Location == "" {
			return strings.Join([]string{s.cfg.BaseURL, filepath}, "/")
		}

		return strings.Join([]string{s.cfg.BaseURL, s.cfg.Location, filepath}, "/")
	}

	return ""
}

// Path returns the filepath prefixed with Location from storage.
func (s *Storage) Path(filepath string) string {
	if _, ok := s.storage.(*fsstorage.Storage); ok {
		return filepath
	}

	return path.Join(s.cfg.Location, filepath)
}

func (s *Storage) Save(ctx context.Context, content io.Reader, path string) error {
	filepath := s.Path(path)
	return s.storage.Save(ctx, content, filepath)
}

func (s *Storage) Stat(ctx context.Context, path string) (*gostorages.Stat, error) {
	return s.storage.Stat(ctx, s.Path(path))
}

func (s *Storage) Open(ctx context.Context, path string) (io.ReadCloser, error) {
	return s.storage.Open(ctx, s.Path(path))
}

func (s *Storage) OpenWithStat(ctx context.Context, path string) (io.ReadCloser, *gostorages.Stat, error) {
	return s.storage.OpenWithStat(ctx, s.Path(path))

}
func (s *Storage) Delete(ctx context.Context, path string) error {
	return s.storage.Delete(ctx, s.Path(path))
}

// New return destination and source storages from config
func New(ctx context.Context, log *slog.Logger, cfg *Config) (*Storage, *Storage, error) {
	if cfg == nil {
		storage := &Storage{storage: &DummyStorage{}}

		log.InfoContext(ctx, "Source storage configured",
			slog.String("type", "dummy"))

		return storage, storage, nil
	}

	var (
		sourceStorage      gostorages.Storage
		destinationStorage gostorages.Storage
		err                error
	)

	if cfg.Source != nil {
		sourceStorage, err = newStorage(ctx, cfg.Source)
		if err != nil {
			return nil, nil, err
		}

		log.InfoContext(ctx, "Source storage configured",
			slog.String("type", cfg.Source.Type))
	}

	if cfg.Destination == nil {
		log.InfoContext(ctx, "Destination storage not set, source storage will be used",
			slog.String("type", cfg.Source.Type))

		return &Storage{
				storage: sourceStorage,
				cfg:     *cfg.Source,
			}, &Storage{
				storage: sourceStorage,
				cfg:     *cfg.Source,
			}, nil
	}

	destinationStorage, err = newStorage(ctx, cfg.Destination)
	if err != nil {
		return nil, nil, err
	}

	log.InfoContext(ctx, "Destination storage configured",
		slog.String("type", cfg.Destination.Type))

	return &Storage{
			storage: sourceStorage,
			cfg:     *cfg.Source,
		}, &Storage{
			storage: destinationStorage,
			cfg:     *cfg.Destination,
		}, nil
}

func newStorage(ctx context.Context, cfg *StorageConfig) (gostorages.Storage, error) {
	if cfg == nil {
		return &DummyStorage{}, nil
	}

	if strings.HasPrefix(cfg.Type, httpStoragePrefix) && cfg.BaseURL == "" {
		return nil, fmt.Errorf("HTTP Wrapper cannot be used without setting *base_url* in config")
	}

	switch cfg.Type {
	case httpS3StorageType:
		cfg.Type = s3StorageType

		storage, err := newStorage(ctx, cfg)
		if err != nil {
			return nil, err
		}

		return NewHTTPStorage(storage, http.NewClient()), nil
	case s3StorageType:
		return s3storage.NewStorage(s3storage.Config{
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
			Region:          cfg.Region,
			Bucket:          cfg.BucketName,
			Endpoint:        cfg.Endpoint,
		})
	case httpDOs3StorageType:
		cfg.Type = DOs3StorageType

		storage, err := newStorage(ctx, cfg)
		if err != nil {
			return nil, err
		}

		return NewHTTPStorage(storage, http.NewClient()), nil
	case DOs3StorageType:
		region := cfg.Region
		return s3storage.NewStorage(s3storage.Config{
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
			Region:          region,
			Bucket:          cfg.BucketName,
			Endpoint:        fmt.Sprintf("https://%s.digitaloceanspaces.com", region),
		})
	case gcsStorageType:
		return gcstorage.NewStorage(ctx, cfg.SecretAccessKey, cfg.BucketName)
	case fsStorageType:
		return fsstorage.NewStorage(fsstorage.Config{Root: cfg.Location}), nil
	case httpFSStorageType:
		cfg.Type = fsStorageType

		storage, err := newStorage(ctx, cfg)
		if err != nil {
			return nil, err
		}

		return NewHTTPStorage(storage, http.NewClient()), nil
	}

	return nil, fmt.Errorf("storage %s does not exist", cfg.Type)
}
