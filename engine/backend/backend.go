package backend

import (
	"fmt"

	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/thoas/picfit/image"
)

// MethodNotImplementedError is an error returned if method is not implemented
var MethodNotImplementedError = errors.New("Not implemented")

// Options is the engine options
type Options struct {
	Color    string
	Degree   int
	Format   imaging.Format
	Height   int
	Images   []image.ImageFile
	Position string
	Quality  int
	Stick    string
	Upscale  bool
	Width    int
}

func (o Options) String() string {
	return fmt.Sprintf("width:%d height:%d quality:%d upscale:%t",
		o.Width, o.Height, o.Quality, o.Upscale)
}

// Engine is an interface to define an image engine
type Backend interface {
	Fit(img *image.ImageFile, options *Options) ([]byte, error)
	Flat(background *image.ImageFile, options *Options) ([]byte, error)
	Flip(img *image.ImageFile, options *Options) ([]byte, error)
	Resize(img *image.ImageFile, options *Options) ([]byte, error)
	Rotate(img *image.ImageFile, options *Options) ([]byte, error)
	String() string
	Thumbnail(img *image.ImageFile, options *Options) ([]byte, error)
}
