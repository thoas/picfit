package backend

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thoas/picfit/image"
)

// MethodNotImplementedError is an error returned if method is not implemented
var MethodNotImplementedError = errors.New("Not implemented")

// Options is the engine options
type Options struct {
	Color    string
	Degree   int
	Format   image.Format
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
	Fit(ctx context.Context, img *image.ImageFile, options *Options) ([]byte, error)
	Flat(ctx context.Context, background *image.ImageFile, options *Options) ([]byte, error)
	Flip(ctx context.Context, img *image.ImageFile, options *Options) ([]byte, error)
	Resize(ctx context.Context, img *image.ImageFile, options *Options) ([]byte, error)
	Rotate(ctx context.Context, img *image.ImageFile, options *Options) ([]byte, error)
	String() string
	Thumbnail(ctx context.Context, img *image.ImageFile, options *Options) ([]byte, error)
}
