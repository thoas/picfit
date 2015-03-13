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

var Extractors = map[string]extractors.Extractor{
	"op":   extractors.Operation,
	"fmt":  extractors.Format,
	"url":  extractors.URL,
	"q":    extractors.Quality,
	"path": extractors.Path,
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 not found", http.StatusNotFound)
	})
}

type Options struct {
	Format  string
	Quality int
}

type Request struct {
	*muxer.Request
	Operation  *image.Operation
	Connection gokvstores.KVStoreConnection
	Key        string
	URL        *url.URL
	Filepath   string
	Options    *Options
}

type Handler func(muxer.Response, *Request)

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	con := App.KVStore.Connection()
	defer con.Close()

	request := muxer.NewRequest(req)

	for k, v := range request.Params {
		request.QueryString[k] = v
	}

	res := muxer.NewResponse(w)

	extracted := map[string]interface{}{}

	for key, extractor := range Extractors {
		result, err := extractor(key, request)

		if err != nil {
			App.Logger.Info(err)

			res.BadRequest()
			return
		}

		extracted[key] = result
	}

	sorted := util.SortMapString(request.QueryString)

	valid := App.IsValidSign(sorted)

	delete(sorted, "sig")

	serialized := hash.Serialize(sorted)

	key := hash.Tokey(serialized)

	App.Logger.Infof("Generating key %s from request: %s", key, serialized)

	var u *url.URL
	var path string
	var format string
	var quality int

	value, ok := extracted["url"]

	if ok && value != nil {
		u = value.(*url.URL)
	}

	value, ok = extracted["path"]

	if ok {
		path = value.(string)
	}

	if !valid || (u == nil && path == "") {
		res.BadRequest()
		return
	}

	value, ok = extracted["fmt"]

	if ok {
		format = value.(string)
	}

	value, ok = extracted["q"]

	if ok && value != nil {
		quality = value.(int)
	}

	options := &Options{Quality: quality, Format: format}

	h(res, &Request{
		request,
		extracted["op"].(*image.Operation),
		con,
		key,
		u,
		path,
		options,
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
