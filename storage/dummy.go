package storage

import (
	"time"

	"github.com/thoas/gostorages"
)

type DummyStorage struct {
}

func (s *DummyStorage) Save(filepath string, file gostorages.File) error {
	return nil
}

func (s *DummyStorage) Path(filepath string) string {
	return ""
}

func (s *DummyStorage) Exists(filepath string) bool {
	return false
}

func (s *DummyStorage) Delete(filepath string) error {
	return nil
}

func (s *DummyStorage) Open(filepath string) (gostorages.File, error) {
	return nil, nil
}

func (s *DummyStorage) ModifiedTime(filepath string) (time.Time, error) {
	return time.Time{}, nil
}

func (s *DummyStorage) Size(filepath string) int64 {
	return 0
}

func (s *DummyStorage) URL(filename string) string {
	return ""
}

func (s *DummyStorage) HasBaseURL() bool {
	return false
}
