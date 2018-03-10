package gostorages

import (
	"bytes"
	"io/ioutil"
	"path"
	"strings"
)

type BaseStorage struct {
	BaseURL  string
	Location string
}

type ContentFile struct {
	*bytes.Reader
}

func (f *ContentFile) Close() error {
	return nil
}

func (f *ContentFile) Size() int64 {
	return int64(f.Len())
}

func (f *ContentFile) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(f)
}

func NewContentFile(content []byte) *ContentFile {
	return &ContentFile{bytes.NewReader(content)}
}

func NewBaseStorage(location string, baseURL string) *BaseStorage {
	return &BaseStorage{
		BaseURL:  strings.TrimSuffix(baseURL, "/"),
		Location: location,
	}
}

func (s *BaseStorage) URL(filename string) string {
	if s.HasBaseURL() {
		return strings.Join([]string{s.BaseURL, s.Path(filename)}, "/")
	}

	return ""
}

func (s *BaseStorage) HasBaseURL() bool {
	return s.BaseURL != ""
}

// Path joins the given file to the storage path
func (s *BaseStorage) Path(filepath string) string {
	return path.Join(s.Location, filepath)
}
