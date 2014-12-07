package application

import (
	"github.com/thoas/kvstores"
	"github.com/thoas/muxer"
	"github.com/thoas/picfit/image"
	"net/http"
	"net/url"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 not found", http.StatusNotFound)
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(NotFound)
}

type Request struct {
	*muxer.Request
	Method     *image.Method
	Connection kvstores.KVStoreConnection
	Key        string
	URL        *url.URL
	Filename   string
}

type Handler func(muxer.Response, *Request)

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	con := App.KVStore.Connection()
	defer con.Close()

	request := muxer.NewRequest(req)

	method, err := extractMethod(request)

	res := muxer.NewResponse(w)

	if err != nil {
		res.BadRequest()
		return
	}

	url, err := extractURL(request)

	filename := request.Params["filename"]

	if err != nil && filename == "" {
		res.BadRequest()
		return
	}

	h(res, &Request{request, method, con, keyFromRequest(request), url, filename})
}

var ImageHandler Handler = func(res muxer.Response, req *Request) {
	imageResponse, err := App.ImageResponseFromRequest(req, true)

	panicIf(err)

	content, err := imageResponse.ToBytes()

	panicIf(err)

	res.ContentType(imageResponse.ContentType)
	res.SetHeaders(imageResponse.Header, true)
	res.ResponseWriter.Write(content)
}
