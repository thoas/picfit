package image

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/franela/goreq"
	"mime"
	"net/http"
)

func ImageFileFromURL(url string) (*ImageFile, error) {
	content, err := goreq.Request{Uri: url}.Do()

	if err != nil {
		return nil, err
	}

	if content.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s [status: %d]", url, content.StatusCode)
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

	var contentType = mime.TypeByExtension(url)

	if value, ok := headers["Content-Type"]; ok {
		contentType = value
	}

	return &ImageFile{Source: dest, ContentType: contentType, Header: headers}, nil
}

func ImageFileFromBytes(content []byte, contentType string, headers map[string]string) (*ImageFile, error) {
	reader := bytes.NewReader(content)

	dest, err := imaging.Decode(reader)

	if err != nil {
		return nil, err
	}

	return &ImageFile{Source: dest, ContentType: contentType, Header: headers}, nil
}
