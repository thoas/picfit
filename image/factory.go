package image

import (
	"io/ioutil"
	"net/url"

	"github.com/ulule/gostorages"

	"github.com/thoas/picfit/storage"
)

// FromURL retrieves an ImageFile from an url
func FromURL(u *url.URL, userAgent string) (*ImageFile, error) {
	storage := &storage.HTTPStorage{UserAgent: userAgent}

	content, err := storage.OpenFromURL(u)
	if err != nil {
		return nil, err
	}

	headers, err := storage.HeadersFromURL(u)
	if err != nil {
		return nil, err
	}

	return &ImageFile{
		Source:   content,
		Headers:  headers,
		Filepath: u.Path[1:],
	}, nil
}

// FromStorage retrieves an ImageFile from storage
func FromStorage(storage gostorages.Storage, filepath string) (*ImageFile, error) {
	var file *ImageFile
	var err error

	f, err := storage.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	modifiedTime, err := storage.ModifiedTime(filepath)
	if err != nil {
		return nil, err
	}

	file = &ImageFile{
		Filepath: filepath,
		Storage:  storage,
	}

	contentType := file.ContentType()

	headers := map[string]string{
		"Last-Modified": modifiedTime.Format(gostorages.LastModifiedFormat),
		"Content-Type":  contentType,
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	file.Source = buf
	file.Headers = headers

	return file, err
}
