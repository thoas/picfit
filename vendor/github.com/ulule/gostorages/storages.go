package gostorages

import (
	"time"
)

type Storage interface {
	Save(filepath string, file File) error
	Path(filepath string) string
	Exists(filepath string) bool
	Delete(filepath string) error
	Open(filepath string) (File, error)
	ModifiedTime(filepath string) (time.Time, error)
	Size(filepath string) int64
	URL(filename string) string
	HasBaseURL() bool
	IsNotExist(err error) bool
}

type File interface {
	Size() int64
	Read(b []byte) (int, error)
	ReadAll() ([]byte, error)
	Close() error
}
