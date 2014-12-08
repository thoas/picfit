package image

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/franela/goreq"
	"net/http"
	"net/url"
)

func ImageFileFromURL(u *url.URL) (*ImageFile, error) {
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

	return &ImageFile{Source: dest, Header: headers, Filepath: u.Path[1:]}, nil
}
