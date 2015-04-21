package application

import (
	"bytes"
	"github.com/mholt/binding"
	"github.com/thoas/gostorages"
	"io"
	"mime/multipart"
)

type MultipartForm struct {
	Data *multipart.FileHeader `json:"data"`
}

func (f *MultipartForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&f.Data: "data",
	}
}

func (f *MultipartForm) Upload(storage gostorages.Storage) error {
	var fh io.ReadCloser

	fh, err := f.Data.Open()

	if err != nil {
		return err
	}

	defer fh.Close()

	dataBytes := bytes.Buffer{}

	_, err = dataBytes.ReadFrom(fh)

	if err != nil {
		return err
	}

	return storage.Save(f.Data.Filename, gostorages.NewContentFile(dataBytes.Bytes()))
}
