package image

import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/thoas/picfit/http"
	"github.com/thoas/storages"
	"net/url"
)

func FromURL(u *url.URL) (*ImageFile, error) {
	storage := &http.HTTPStorage{}

	content, err := storage.OpenFromURL(u)

	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(content)

	dest, err := imaging.Decode(reader)

	if err != nil {
		return nil, err
	}

	headers, err := storage.HeadersFromURL(u)

	if err != nil {
		return nil, err
	}

	return &ImageFile{
		Source:   dest,
		Headers:  headers,
		Filepath: u.Path[1:],
	}, nil
}

func FromStorage(storage storages.Storage, filepath string) (*ImageFile, error) {
	var file *ImageFile
	var err error

	body, err := storage.Open(filepath)

	if err != nil {
		return nil, err
	}

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
		"Last-Modified": modifiedTime.Format(storages.LastModifiedFormat),
		"Content-Type":  contentType,
	}

	reader := bytes.NewReader(body)

	dest, err := imaging.Decode(reader)

	if err != nil {
		return nil, err
	}

	file.Source = dest
	file.Headers = headers

	return file, err
}
