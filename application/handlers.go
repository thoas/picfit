package application

import (
	"encoding/json"
	"github.com/thoas/gokvstores"
	"github.com/thoas/muxer"
	"github.com/thoas/picfit/extractors"
	"github.com/thoas/picfit/hash"
	"github.com/thoas/picfit/image"
	"github.com/thoas/picfit/util"
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
	Operation  *image.Operation
	Connection gokvstores.KVStoreConnection
	Key        string
	URL        *url.URL
	Filepath   string
	Format     string
}

type Handler func(muxer.Response, *Request)

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	con := App.KVStore.Connection()
	defer con.Close()

	request := muxer.NewRequest(req)

	for k, v := range request.Params {
		request.QueryString[k] = v
	}

	operation, errop := extractors.Operation(request)

	res := muxer.NewResponse(w)

	url, err := extractors.URL(request)

	filepath, ok := request.QueryString["path"]

	format, errfmt := extractors.Format(request)

	sorted := util.SortMapString(request.QueryString)

	valid := App.IsValidSign(sorted)

	delete(sorted, "sig")

	serialized := hash.Serialize(sorted)

	key := hash.Tokey(serialized)

	App.Logger.Infof("Generating key %s from request: %s", key, serialized)

	if (err != nil && !ok) || errfmt != nil || errop != nil || !valid {
		res.BadRequest()
		return
	}

	h(res, &Request{
		request,
		operation,
		con,
		key,
		url,
		filepath,
		format,
	})
}

var ImageHandler Handler = func(res muxer.Response, req *Request) {
	file, err := App.ImageFileFromRequest(req, true, true)

	util.PanicIf(err)

	content, err := file.ToBytes()

	util.PanicIf(err)

	res.SetHeaders(file.Headers, true)
	res.ResponseWriter.Write(content)
}

var GetHandler Handler = func(res muxer.Response, req *Request) {
	file, err := App.ImageFileFromRequest(req, false, false)

	util.PanicIf(err)

	content, err := json.Marshal(map[string]string{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})

	util.PanicIf(err)

	res.ContentType("application/json")
	res.ResponseWriter.Write(content)
}

var RedirectHandler Handler = func(res muxer.Response, req *Request) {
	file, err := App.ImageFileFromRequest(req, false, false)

	util.PanicIf(err)

	res.PermanentRedirect(file.URL())
}
