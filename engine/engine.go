package engine

import (
	"github.com/disintegration/imaging"

	"github.com/thoas/picfit/image"
)

type Options struct {
	Upscale bool
	Format  imaging.Format
	Quality int
}

type Engine interface {
	Resize(img *image.ImageFile, width int, height int, options *Options) ([]byte, error)
	Thumbnail(img *image.ImageFile, width int, height int, options *Options) ([]byte, error)
	Transform(img *image.ImageFile, operation *Operation, qs map[string]string) (*image.ImageFile, error)
	Flip(img *image.ImageFile, pos string, options *Options) ([]byte, error)
	Rotate(img *image.ImageFile, deg int, options *Options) ([]byte, error)
	Fit(img *image.ImageFile, width int, height int, options *Options) ([]byte, error)
}
