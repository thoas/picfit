package extractors

import (
	"fmt"
	"github.com/thoas/muxer"
	"github.com/thoas/picfit/image"
	"mime"
	"net/url"
	"path/filepath"
)

type Extractor func(key string, req *muxer.Request) (interface{}, error)

var Operation Extractor = func(key string, req *muxer.Request) (interface{}, error) {
	operation, ok := image.Operations[req.QueryString[key]]

	if !ok {
		return nil, fmt.Errorf("Invalid method %s or invalid parameters", operation)
	}

	return operation, nil
}

var URL Extractor = func(key string, req *muxer.Request) (interface{}, error) {
	value, ok := req.QueryString[key]

	if !ok {
		return nil, nil
	}

	url, err := url.Parse(value)

	if err != nil {
		return nil, fmt.Errorf("URL %s is not valid", value)
	}

	mimetype := mime.TypeByExtension(filepath.Ext(value))

	_, ok = image.Extensions[mimetype]

	if !ok {
		return nil, fmt.Errorf("Mimetype %s is not supported", mimetype)
	}

	return url, nil
}

var Path Extractor = func(key string, req *muxer.Request) (interface{}, error) {
	return req.QueryString[key], nil
}
