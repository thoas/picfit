package fastimage

// ImageType represents the type of the image detected, or `Unknown`.
type ImageType uint

//go:generate stringer -type=ImageType -output=image_type_string.go
const (
	// GIF represents a GIF image
	GIF ImageType = iota
	// PNG represents a PNG image
	PNG
	// JPEG represents a JPEG image
	JPEG
	// BMP represents a BMP image
	BMP
	// TIFF represents a TIFF image
	TIFF
	// Unknown represents an unknown image type
	Unknown
)
