package image

import (
	"bytes"
	"context"
	"io"
	"net/url"

	"github.com/thoas/picfit/constants"
	"github.com/thoas/picfit/http"
	storagepkg "github.com/thoas/picfit/storage"
)

// FromURL retrieves an ImageFile from an url
func FromURL(u *url.URL, userAgent string) (*ImageFile, error) {
	storage := storagepkg.NewHTTPStorage(nil, http.NewClient(http.WithUserAgent(userAgent)))

	content, err := storage.OpenFromURL(u)
	if err != nil {
		return nil, err
	}

	headers, err := storage.HeadersFromURL(u)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, content); err != nil {
		return nil, err
	}
	if err := content.Close(); err != nil {
		return nil, err
	}
	return &ImageFile{
		Source:   buffer.Bytes(),
		Headers:  headers,
		Filepath: u.Path[1:],
	}, nil
}

// FromStorage retrieves an ImageFile from storage
func FromStorage(ctx context.Context, storage storagepkg.Storage, filepath string) (*ImageFile, error) {
	var file *ImageFile
	var err error

	f, err := storage.Open(ctx, filepath)
	if err != nil {
		return nil, err
	}

	stat, err := storage.Stat(ctx, filepath)
	if err != nil {
		return nil, err
	}

	file = &ImageFile{
		Filepath: filepath,
		Storage:  storage,
	}

	contentType := file.ContentType()

	headers := map[string]string{
		"Last-Modified": stat.ModifiedTime.Format(constants.ModifiedTimeFormat),
		"Content-Type":  contentType,
	}

	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, f); err != nil {
		return nil, err
	}

	file.Source = buffer.Bytes()
	file.Headers = headers
	if err := f.Close(); err != nil {
		return nil, err
	}
	return file, err
}
