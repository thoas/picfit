package backend

import (
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/thoas/picfit/image"
)

// MethodNotImplementedError is an error returned if method is not implemented
var MethodNotImplementedError = errors.New("Not implemented")

// Options is the engine options
type Options struct {
	Upscale  bool
	Format   imaging.Format
	Quality  int
	Width    int
	Height   int
	Position string
	Color    string
	Degree   int
	Images   []image.ImageFile
}

// Engine is an interface to define an image engine
type Backend interface {
	String() string
	Resize(img *image.ImageFile, options *Options) ([]byte, error)
	Thumbnail(img *image.ImageFile, options *Options) ([]byte, error)
	Flip(img *image.ImageFile, options *Options) ([]byte, error)
	Rotate(img *image.ImageFile, options *Options) ([]byte, error)
	Fit(img *image.ImageFile, options *Options) ([]byte, error)
	Flat(background *image.ImageFile, options *Options) ([]byte, error)
}
