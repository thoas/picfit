package backend

import (
	"github.com/thoas/picfit/engine/config"
	imagefile "github.com/thoas/picfit/image"
)

type Vips struct {
}

func NewVips(cfg config.Config) *Vips {
	return &Vips{}
}

func (e *Vips) String() string {
	return "vips"
}

func (e *Vips) Fit(background *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Vips) Flat(background *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Vips) Flip(background *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Vips) Resize(background *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Vips) Rotate(background *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Vips) Thumbnail(background *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}
