package application

import (
	"errors"
	"fmt"
	"github.com/thoas/muxer"
	"github.com/thoas/picfit/image"
	"mime"
	"net/http"
	"net/url"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 not found", http.StatusNotFound)
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(NotFound)
}

func extractMethod(req *muxer.Request) (*image.Method, error) {
	method, ok := image.Methods[req.Params["method"]]

	if !ok || !method.HasEnoughParams(req.QueryString) {
		return nil, errors.New(fmt.Sprintf("Invalid method %s or invalid parameters", method))
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

	return nil, errors.New(fmt.Sprintf("URL %s is not valid", urlValue))
}

func keyFromRequest(req *muxer.Request) string {
	qs := serialize(req.QueryString)

	var key string

	if filename, ok := req.Params["filename"]; !ok {
		key = tokey(req.Params["method"], qs)
	} else {
		key = tokey(req.Params["method"], filename, qs)
	}

	return key
}

var ImageHandler muxer.Handler = func(res muxer.Response, req *muxer.Request) {
	method, err := extractMethod(req)

	panicIf(err)

	url, err := extractURL(req)

	filename := req.Params["filename"]

	if err != nil && filename == "" {
		res.BadRequest()
		return
	}

	con := App.KVStore.Connection()
	defer con.Close()

	key := keyFromRequest(req)

	stored := con.Get(key)

	var imageResponse *image.ImageResponse

	// Image from the KVStore found
	if stored != "" {
		// URL provided we use http protocol to retrieve it
		if App.BaseURL != "" {
			imageResponse, err = image.ImageResponseFromURL(App.URL(stored))

			panicIf(err)
		} else {
			imageResponse, err = App.ImageResponseFromStorage(stored)

			panicIf(err)
		}
	} else {
		// Image not found from the KVStore, we need to process it
		// URL available in Query String
		if url != nil {
			imageResponse, err = image.ImageResponseFromURL(url.String())

			panicIf(err)
		} else {
			// URL provided we use http protocol to retrieve it
			if App.BaseURL != "" {
				imageResponse, err = image.ImageResponseFromURL(App.URL(filename))

				panicIf(err)
			} else {
				imageResponse, err = App.ImageResponseFromStorage(stored)

				panicIf(err)
			}
		}

		file := image.NewImageFile(imageResponse.Image)

		dest, err := file.Transform(method, req.QueryString)

		panicIf(err)

		imageResponse.Image = dest
		imageResponse.Key = key

		go App.Store(imageResponse)
	}

	content, err := imageResponse.ToBytes()

	panicIf(err)

	res.ResponseWriter.Write(content)
}
