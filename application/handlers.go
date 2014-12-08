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
	Filepath   string
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

	filepath := request.Params["filepath"]

	if err != nil && filepath == "" {
		res.BadRequest()
		return
	}

	h(res, &Request{request, method, con, keyFromRequest(request), url, filepath})
}

var ImageHandler Handler = func(res muxer.Response, req *Request) {
	file, err := App.ImageFileFromRequest(req, true, true)

	panicIf(err)

	content, err := file.ToBytes()

	panicIf(err)

	res.SetHeaders(file.Header, true)
	res.ResponseWriter.Write(content)
}

var GetHandler Handler = func(res muxer.Response, req *Request) {
	file, err := App.ImageFileFromRequest(req, false, false)

	panicIf(err)

	content, err := App.ToJSON(file)

	panicIf(err)

	res.ContentType("application/json")
	res.ResponseWriter.Write(content)
}
