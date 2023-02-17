package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ulule/gostorages"

	"github.com/thoas/picfit/failure"
	httppkg "github.com/thoas/picfit/http"
)

// HTTPStorage wraps a storage
type HTTPStorage struct {
	gostorages.Storage
	httpclient *httppkg.Client
}

// HeaderKeys represents the list of headers
var HeaderKeys = []string{
	"Age",
	"Content-Type",
	"Last-Modified",
	"Date",
	"Etag",
}

func NewHTTPStorage(storage gostorages.Storage, httpclient *httppkg.Client) *HTTPStorage {
	return &HTTPStorage{
		Storage:    storage,
		httpclient: httpclient,
	}
}

// Open retrieves a gostorages File from a filepath
func (s *HTTPStorage) Open(filepath string) (gostorages.File, error) {
	u, err := url.Parse(s.URL(filepath))
	if err != nil {
		return nil, err
	}

	content, err := s.OpenFromURL(u)
	if err != nil {
		return nil, err
	}

	return gostorages.NewContentFile(content), nil
}

// OpenFromURL retrieves bytes from an url
func (s *HTTPStorage) OpenFromURL(u *url.URL) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpclient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, failure.ErrFileNotExists
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s [status: %d]", u.String(), resp.StatusCode)
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// HeadersFromURL retrieves the headers from an url
func (s *HTTPStorage) HeadersFromURL(u *url.URL) (map[string]string, error) {
	var headers = make(map[string]string)

	req, err := http.NewRequestWithContext(context.Background(), "HEAD", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpclient.Do(req)
	if err != nil {
		return nil, err
	}

	for _, key := range HeaderKeys {
		if value, ok := resp.Header[key]; ok && len(value) > 0 {
			headers[key] = value[0]
		}
	}
	return headers, nil
}

// Headers returns headers from a filepath
func (s *HTTPStorage) Headers(filepath string) (map[string]string, error) {
	u, err := url.Parse(s.URL(filepath))
	if err != nil {
		return nil, err
	}

	return s.HeadersFromURL(u)
}

// ModifiedTime returns the modified time from a filepath
func (s *HTTPStorage) ModifiedTime(filepath string) (time.Time, error) {
	headers, err := s.Headers(filepath)
	if err != nil {
		return time.Time{}, err
	}

	lastModified, ok := headers["Last-Modified"]
	if !ok {
		return time.Time{}, fmt.Errorf("Last-Modified header not found")
	}

	return time.Parse(gostorages.LastModifiedFormat, lastModified)
}

func (s *HTTPStorage) IsNotExist(err error) bool {
	return false
}
