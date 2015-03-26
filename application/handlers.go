package application

import (
	"encoding/json"
	"github.com/thoas/gokvstores"
	"github.com/thoas/muxer"
	"github.com/thoas/picfit/engines"
	"github.com/thoas/picfit/extractors"
	"github.com/thoas/picfit/util"
	"net/http"
	"net/url"
)

var Extractors = map[string]extractors.Extractor{
	"op":   extractors.Operation,
	"url":  extractors.URL,
	"path": extractors.Path,
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 not found", http.StatusNotFound)
	})
}

type Request struct {
	*muxer.Request
	Operation  *engines.Operation
	Connection gokvstores.KVStoreConnection
	Key        string
	URL        *url.URL
	Filepath   string
}

type Handler func(muxer.Response, *Request, *Application)

var ImageHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, true, true)

	util.PanicIf(err)

	res.SetHeaders(file.Headers, true)
	res.ResponseWriter.Write(file.Content())
}

var GetHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, false, false)

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

var RedirectHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, false, false)

	util.PanicIf(err)

	res.PermanentRedirect(file.URL())
}
