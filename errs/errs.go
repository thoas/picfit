package errs

import (
	"errors"
	"net/http"
)

// ErrFileNotExists is a targetted error when image does not exist on storage
var ErrFileNotExists = errors.New("File does not exist")

// ErrKeyNotExists is a targetted error when image does not exist on storage
var ErrKeyNotExists = errors.New("Key does not exist")

// Handle returns the proper http code based on an error
func Handle(err error, response http.ResponseWriter) {
	if err == ErrFileNotExists || err == ErrKeyNotExists {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	panic(err)
}
