package extractors

import (
	"fmt"
	"github.com/thoas/muxer"
	"github.com/thoas/picfit/image"
	"mime"
	"net/url"
)

func Operation(req *muxer.Request) (*image.Operation, error) {
	operation, ok := image.Operations[req.QueryString["op"]]

	if !ok {
		return nil, fmt.Errorf("Invalid method %s or invalid parameters", operation)
	}

	return operation, nil
}

func URL(req *muxer.Request) (*url.URL, error) {
	value, ok := req.QueryString["url"]

	if ok {
		url, err := url.Parse(value)

		mimetype := mime.TypeByExtension(value)

		_, ok := image.Formats[mimetype]

		if ok || err == nil {
			return url, nil
		}
	}

	return nil, fmt.Errorf("URL %s is not valid", value)
}
