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

	if !ok || !operation.HasEnoughParams(req.QueryString) {
		return nil, fmt.Errorf("Invalid method %s or invalid parameters", operation)
	}

	return operation, nil
}

func URL(req *muxer.Request) (*url.URL, error) {
	urlValue := req.QueryString["url"]

	if urlValue != "" {
		url, err := url.Parse(urlValue)

		mimetype := mime.TypeByExtension(urlValue)

		_, ok := image.Formats[mimetype]

		if ok || err == nil {
			return url, nil
		}
	}

	return nil, fmt.Errorf("URL %s is not valid", urlValue)
}
