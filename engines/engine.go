package engines

import (
	"github.com/thoas/picfit/image"
)

type Options struct {
	Upscale bool
	Format  string
	Quality int
}

type Engine interface {
	Resize(img *image.ImageFile, width int, height int, options *Options) ([]byte, error)
	Thumbnail(img *image.ImageFile, width int, height int, options *Options) ([]byte, error)
	Transform(img *image.ImageFile, operation *Operation, qs map[string]string) (*image.ImageFile, error)
}
