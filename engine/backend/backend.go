package backend

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/thoas/picfit/image"
)

// MethodNotImplementedError is an error returned if method is not implemented
var MethodNotImplementedError = errors.New("Not implemented")

// Options is the engine options
type Options struct {
	Color    string
	Degree   int
	Filter   string
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
	Effect(ctx context.Context, dst io.Writer, img *image.ImageFile, options *Options) error
	Fit(ctx context.Context, dst io.Writer, img *image.ImageFile, options *Options) error
	Flat(ctx context.Context, dst io.Writer, background *image.ImageFile, options *Options) error
	Flip(ctx context.Context, dst io.Writer, img *image.ImageFile, options *Options) error
	Resize(ctx context.Context, dst io.Writer, img *image.ImageFile, options *Options) error
	Rotate(ctx context.Context, dst io.Writer, img *image.ImageFile, options *Options) error
	String() string
	Thumbnail(ctx context.Context, dst io.Writer, img *image.ImageFile, options *Options) error
}
