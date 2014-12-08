package application

import (
	"fmt"
	"github.com/thoas/muxer"
	"github.com/thoas/picfit/image"
	"mime"
	"net/url"
)

func extractMethod(req *muxer.Request) (*image.Method, error) {
	method, ok := image.Methods[req.Params["method"]]

	if !ok || !method.HasEnoughParams(req.QueryString) {
		return nil, fmt.Errorf("Invalid method %s or invalid parameters", method)
	}

	return method, nil
}

func extractURL(req *muxer.Request) (*url.URL, error) {
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
