package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mholt/binding"
	"github.com/thoas/gostorages"
	"github.com/thoas/muxer"
	"io"
	"mime/multipart"
	"net/http"
)

type MultipartForm struct {
	Data *multipart.FileHeader `json:"data"`
}

func (f *MultipartForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&f.Data: "data",
	}
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 not found", http.StatusNotFound)
	})
}

type Handler func(muxer.Response, *Request, *Application)

var ImageHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, true, true)

	if err != nil {
		panic(err)
	}

	res.SetHeaders(file.Headers, true)
	res.ResponseWriter.Write(file.Content())
}

var UploadHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	if app.SourceStorage == nil {
		res.Abort(500, "Your application doesn't have a source storage")
		return
	}

	multipartForm := new(MultipartForm)
	errs := binding.Bind(req.Request.Request, multipartForm)
	if errs.Handle(res) {
		return
	}

	var fh io.ReadCloser
	var err error
	if fh, err = multipartForm.Data.Open(); err != nil {
		res.Abort(500, fmt.Sprint("Error opening Mime::Data %+v", err))
		return
	}
	defer fh.Close()
	dataBytes := bytes.Buffer{}
	var size int64
	if size, err = dataBytes.ReadFrom(fh); err != nil {
		res.Abort(500, fmt.Sprint("Error reading Mime::Data %+v", err))
		return
	}

	app.Logger.Printf("Read %v bytes with filename %s", size, multipartForm.Data.Filename)

	err = app.SourceStorage.Save(multipartForm.Data.Filename, gostorages.NewContentFile(dataBytes.Bytes()))

	if err != nil {
		res.Abort(500, fmt.Sprint("Error uploading file %s to source storage %+v", multipartForm.Data.Filename, err))
		return
	}
}

var GetHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, false, false)

	if err != nil {
		panic(err)
	}

	content, err := json.Marshal(map[string]string{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})

	if err != nil {
		panic(err)
	}

	res.ContentType("application/json")
	res.ResponseWriter.Write(content)
}

var RedirectHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, false, false)

	if err != nil {
		panic(err)
	}

	res.PermanentRedirect(file.URL())
}
