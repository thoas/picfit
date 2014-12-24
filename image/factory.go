package image

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/franela/goreq"
	"github.com/thoas/storages"
	"net/http"
	"net/url"
)

func FromURL(u *url.URL) (*ImageFile, error) {
	content, err := goreq.Request{Uri: u.String()}.Do()

	if err != nil {
		return nil, err
	}

	if content.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s [status: %d]", u.String(), content.StatusCode)
	}

	dest, err := imaging.Decode(content.Body)

	if err != nil {
		return nil, err
	}

	var headers = make(map[string]string)

	for _, key := range HeaderKeys {
		if value, ok := content.Header[key]; ok && len(value) > 0 {
			headers[key] = value[0]
		}
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

	// URL provided we use http protocol to retrieve it
	if storage.HasBaseURL() {
		u, err := url.Parse(storage.URL(filepath))

		if err != nil {
			return nil, err
		}

		file, err = FromURL(u)

		if err != nil {
			return nil, err
		}
	} else {
		body, err := storage.Open(filepath)

		if err != nil {
			return nil, err
		}

		modifiedTime, err := storage.ModifiedTime(filepath)

		if err != nil {
			return nil, err
		}

		i := &ImageFile{Filepath: filepath}

		contentType := i.ContentType()

		headers := map[string]string{
			"Last-Modified": modifiedTime.Format(storages.LastModifiedFormat),
			"Content-Type":  contentType,
		}

		reader := bytes.NewReader(body)

		dest, err := imaging.Decode(reader)

		if err != nil {
			return nil, err
		}

		return &ImageFile{
			Source:   dest,
			Storage:  storage,
			Headers:  headers,
			Filepath: filepath,
		}, nil
	}

	return file, err
}
