package payload

import (
	"mime/multipart"
	"net/http"

	"github.com/mholt/binding"
)

// Multipart represents a multipart upload
type Multipart struct {
	Data *multipart.FileHeader `json:"data"`
}

// FieldMap defines excepted inputs
func (f *Multipart) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&f.Data: "data",
	}
}
