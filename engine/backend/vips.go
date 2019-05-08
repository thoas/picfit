package backend

import (
	"github.com/thoas/picfit/engine/config"
	imagefile "github.com/thoas/picfit/image"
	//"gopkg.in/h2non/bimg.v1"
)

type Vips struct {
}

func NewVips(cfg config.Config) *Vips {
	return &Vips{}
}

func (e *Vips) String() string {
	return "vips"
}

func (e *Vips) Fit(image *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Vips) Fill(image *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Vips) Flat(image *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Vips) Flip(image *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Vips) Resize(image *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Vips) Rotate(image *imagefile.ImageFile, options *Options) ([]byte, error) {
	return nil, MethodNotImplementedError
}

func (e *Vips) Thumbnail(image *imagefile.ImageFile, options *Options) ([]byte, error) {
	/*
		bimgOptions := bimg.Options{
			Width:   options.Width,
			Height:  options.Height,
			Gravity: bimg.GravitySmart,
			//Gravity:   bimg.GravityNorth,
			Crop:      true,
		}
		return bimg.NewImage(image.Source).Process(bimgOptions)
	*/
	return nil, MethodNotImplementedError
}
