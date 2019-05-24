package failure

import (
	"errors"
)

var (
	// ErrFileNotExists is an error when image does not exist on storage
	ErrFileNotExists = errors.New("File does not exist")

	// ErrKeyNotExists is an error when image does not exist on storage
	ErrKeyNotExists = errors.New("Key does not exist")

	// ErrQuality is an error when the quality requested is higher than expected
	ErrQuality = errors.New("Quality should be <= 100")

	// ErrUnprocessable is an error when parameters are missing
	ErrUnprocessable = errors.New("Unprocessable request, missing parameters")

	// ErrFileNotModified is an error when file is not modified
	ErrFileNotModified = errors.New("File not modified")
)
