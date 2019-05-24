package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/franela/goreq"

	"github.com/ulule/gostorages"

	"github.com/thoas/picfit/failure"
)

// HTTPStorage wraps a storage
type HTTPStorage struct {
	gostorages.Storage
	UserAgent string
}

// HeaderKeys represents the list of headers
var HeaderKeys = []string{
	"Age",
	"Content-Type",
	"Last-Modified",
	"Date",
	"Etag",
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
	content, err := goreq.Request{
		Uri:       u.String(),
		UserAgent: s.UserAgent,
	}.Do()
	if err != nil {
		return nil, err
	}

	defer content.Body.Close()

	if content.StatusCode == http.StatusNotFound {
		return nil, failure.ErrFileNotExists
	}

	if content.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s [status: %d]", u.String(), content.StatusCode)
	}

	return ioutil.ReadAll(content.Body)
}

// HeadersFromURL retrieves the headers from an url
func (s *HTTPStorage) HeadersFromURL(u *url.URL) (map[string]string, error) {
	var headers = make(map[string]string)

	content, err := goreq.Request{
		Uri:       u.String(),
		Method:    "GET",
		UserAgent: s.UserAgent,
	}.Do()
	if err != nil {
		return nil, err
	}

	defer content.Body.Close()

	for _, key := range HeaderKeys {
		if value, ok := content.Header[key]; ok && len(value) > 0 {
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
