package application

import (
	"bytes"
	"github.com/mholt/binding"
	"github.com/thoas/gostorages"
	"github.com/thoas/picfit/image"
	"io"
	"mime/multipart"
	"net/http"
)

type MultipartForm struct {
	Data *multipart.FileHeader `json:"data"`
}

func (f *MultipartForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&f.Data: "data",
	}
}

func (f *MultipartForm) Upload(storage gostorages.Storage) (*image.ImageFile, error) {
	var fh io.ReadCloser

	fh, err := f.Data.Open()

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	dataBytes := bytes.Buffer{}

	_, err = dataBytes.ReadFrom(fh)

	if err != nil {
		return nil, err
	}

	err = storage.Save(f.Data.Filename, gostorages.NewContentFile(dataBytes.Bytes()))

	if err != nil {
		return nil, err
	}

	return &image.ImageFile{
		Filepath: f.Data.Filename,
		Storage:  storage,
	}, nil
}
