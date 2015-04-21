package application

import (
	"github.com/mholt/binding"
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
