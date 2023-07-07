package storage

import (
	"context"
	"io"

	"github.com/ulule/gostorages"
)

type DummyStorage struct {
}

func (s *DummyStorage) Save(ctx context.Context, content io.Reader, filepath string) error {
	return nil
}

func (s *DummyStorage) Delete(ctx context.Context, filepath string) error {
	return nil
}

func (s *DummyStorage) Open(ctx context.Context, filepath string) (io.ReadCloser, error) {
	return nil, nil
}

func (s DummyStorage) Stat(ctx context.Context, path string) (*gostorages.Stat, error) {
	return &gostorages.Stat{}, nil
}

func (s DummyStorage) OpenWithStat(ctx context.Context, path string) (io.ReadCloser, *gostorages.Stat, error) {
	return nil, &gostorages.Stat{}, nil
}
