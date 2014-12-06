package image

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/franela/goreq"
	"image"
	"mime"
	"net/http"
)

type ImageResponse struct {
	Image       image.Image
	ContentType string
	Key         string
}

func ImageResponseFromURL(url string) (*ImageResponse, error) {
	content, err := goreq.Request{Uri: url}.Do()

	if err != nil {
		return nil, err
	}

	if content.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprint("%s [status: %d]", url, content.StatusCode))
	}

	dest, err := imaging.Decode(content.Body)

	if err != nil {
		return nil, err
	}

	var contentType = mime.TypeByExtension(url)

	if results, ok := content.Header["Content-Type"]; ok {
		contentType = results[0]
	}

	return &ImageResponse{Image: dest, ContentType: contentType}, nil
}

func ImageResponseFromBytes(content []byte, contentType string) (*ImageResponse, error) {
	reader := bytes.NewReader(content)

	dest, err := imaging.Decode(reader)

	if err != nil {
		return nil, err
	}

	return &ImageResponse{Image: dest, ContentType: contentType}, nil
}

func (i *ImageResponse) ToBytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := imaging.Encode(buf, i.Image, Formats[i.ContentType])

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (i *ImageResponse) Format() string {
	return Extensions[i.ContentType]
}
